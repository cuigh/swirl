package scaler

import (
	"fmt"
	"strings"
	"time"

	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/data/set"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/auxo/util/run"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/docker"
	"github.com/cuigh/swirl/misc"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

type scaleType int

const (
	scaleNone scaleType = iota
	scaleUp
	scaleDown
)

type policyType string

const (
	policyAny policyType = "any"
	policyAll policyType = "all"
)

const Name = "scaler"

type Scaler struct {
	d        *docker.Docker
	checkers map[string]Checker
}

func NewScaler(d *docker.Docker, mb biz.MetricBiz) *Scaler {
	s := &Scaler{
		d:        d,
		checkers: map[string]Checker{},
	}
	if mb.Enabled() {
		s.checkers["cpu"] = &cpuChecker{mb: mb}
	}
	return s
}

// Start starts a timer to scale services automatically.
func (s *Scaler) Start() {
	const labelScale = "swirl.scale"

	if len(s.checkers) > 0 {
		return
	}

	run.Schedule(time.Minute, func() {
		ctx, cancel := misc.Context(time.Minute)
		defer cancel()

		args := filters.NewArgs()
		args.Add("mode", "replicated")
		args.Add("label", labelScale)
		services, err := s.d.ServiceSearch(ctx, args)
		if err != nil {
			log.Get("scaler").Error("scaler > Failed to search service: ", err)
			return
		}

		for _, service := range services {
			label := service.Spec.Labels[labelScale]
			opts := data.ParseOptions(label, ",", "=")
			s.tryScale(&service, opts)
		}
	}, nil)
}

// nolint: gocyclo
func (s *Scaler) tryScale(service *swarm.Service, opts data.Options) {
	// only care about services running in replicated mode
	if service.Spec.Mode.Replicated == nil {
		return
	}

	var (
		min    uint64 = 2
		max    uint64 = 8
		step   uint64 = 2
		window        = 3 * time.Minute
		policy        = policyAny
		args   data.Options
	)

	for _, opt := range opts {
		switch opt.Name {
		case "min":
			min = cast.ToUint64(opt.Value, min)
		case "max":
			max = cast.ToUint64(opt.Value, max)
		case "step":
			step = cast.ToUint64(opt.Value, step)
		case "window":
			window = cast.ToDuration(opt.Value, window)
		case "policy":
			policy = policyType(opt.Value)
		default:
			args = append(args, opt)
		}
	}

	// ignore services that have been updated in window period.
	if service.UpdatedAt.Add(window).After(time.Now()) {
		return
	}

	result := s.check(service, policy, args)
	if result.Type == scaleNone {
		return
	}

	ctx, cancel := misc.Context(time.Minute)
	defer cancel()

	logger := log.Get("scaler")
	replicas := *service.Spec.Mode.Replicated.Replicas
	if result.Type == scaleUp {
		if replicas < max {
			if err := s.d.ServiceScale(ctx, service.Spec.Name, service.Version.Index, replicas+step); err != nil {
				logger.Errorf("scaler > Failed to scale service '%s': %v", service.Spec.Name, err)
			} else {
				logger.Infof("scaler > Service '%s' scaled up for: %v", service.Spec.Name, result.Reasons)
			}
		}
	} else if result.Type == scaleDown {
		if replicas > min {
			if err := s.d.ServiceScale(ctx, service.Spec.Name, service.Version.Index, replicas-step); err != nil {
				logger.Errorf("scaler > Failed to scale service '%s': %v", service.Spec.Name, err)
			} else {
				logger.Infof("scaler > Service '%s' scaled down for: %v", service.Spec.Name, result.Reasons)
			}
		}
	}
}

func (s *Scaler) check(service *swarm.Service, policy policyType, args data.Options) checkResult {
	result := checkResult{
		Reasons: make(map[string]float64),
	}
	if policy == policyAny {
		for _, arg := range args {
			st, value := s.checkArg(service.Spec.Name, arg)
			if st == scaleNone {
				continue
			}
			result.Type = st
			result.Reasons[arg.Name] = value
			break
		}
	} else if policy == policyAll {
		types := set.Set{}
		for _, arg := range args {
			st, value := s.checkArg(service.Spec.Name, arg)
			types.Add(st)
			if types.Len() > 1 {
				result.Type = scaleNone
				return result
			}
			result.Type = st
			result.Reasons[arg.Name] = value
		}
	}
	return result
}

func (s *Scaler) checkArg(service string, arg data.Option) (scaleType, float64) {
	items := strings.Split(arg.Value, ":")
	if len(items) != 2 {
		log.Get("scaler").Warnf("scaler > Invalid scale argument: %s=%s", arg.Name, arg.Value)
		return scaleNone, 0
	}

	c := s.checkers[arg.Name]
	if c == nil {
		log.Get("scaler").Warnf("scaler > Metric checker '%s' not found", arg.Name)
		return scaleNone, 0
	}

	low := cast.ToFloat64(items[0])
	high := cast.ToFloat64(items[1])
	return c.Check(service, low, high)
}

type checkResult struct {
	Type    scaleType
	Reasons map[string]float64
}

type Checker interface {
	Check(service string, low, high float64) (scaleType, float64)
}

type cpuChecker struct {
	mb biz.MetricBiz
}

func (c *cpuChecker) Check(service string, low, high float64) (scaleType, float64) {
	ctx, cancel := misc.Context(time.Minute)
	defer cancel()

	query := fmt.Sprintf(`avg(rate(container_cpu_user_seconds_total{container_label_com_docker_swarm_service_name="%s"}[1m]) * 100)`, service)
	vector, err := c.mb.GetVector(ctx, query, "", time.Now())
	if err != nil {
		log.Get("scaler").Error("scaler > Failed to query metrics: ", err)
		return scaleNone, 0
	}
	if len(vector.Data) == 0 {
		return scaleNone, 0
	}

	cv := vector.Data[0]
	if cv.Value <= low {
		return scaleDown, cv.Value
	} else if cv.Value >= high {
		return scaleUp, cv.Value
	}
	return scaleNone, 0
}

func Start() error {
	s, err := container.TryFind(Name)
	if err == nil {
		s.(*Scaler).Start()
	}
	return err
}

func init() {
	container.Put(NewScaler, container.Name(Name))
}
