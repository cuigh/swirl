package scaler

import (
	"fmt"
	"strings"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/data/set"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/auxo/util/run"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

type checker func(service string, low, high float64) (scaleType, float64)

var checkers = map[string]checker{
	"cpu": cpuChecker,
}

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

// Start starts a timer to scale services automatically.
func Start() {
	const labelScale = "swirl.scale"

	run.Schedule(time.Minute, func() {
		args := filters.NewArgs()
		args.Add("mode", "replicated")
		args.Add("label", labelScale)
		services, err := docker.ServiceSearch(args)
		if err != nil {
			log.Get("scaler").Error("scaler > Failed to search service: ", err)
			return
		}

		for _, service := range services {
			label := service.Spec.Labels[labelScale]
			opts := data.ParseOptions(label, ",", "=")
			tryScale(&service, opts)
		}
	}, nil)
}

// nolint: gocyclo
func tryScale(service *swarm.Service, opts data.Options) {
	// ignore services with global mode
	if service.Spec.Mode.Replicated == nil {
		return
	}

	// ignore services that have been updated in 3 minutes
	if service.UpdatedAt.Add(3 * time.Minute).After(time.Now()) {
		return
	}

	var (
		min    = uint64(2)
		max    = uint64(5)
		policy = policyAny
		args   data.Options
	)
	for _, opt := range opts {
		switch opt.Name {
		case "min":
			min = cast.ToUint64(opt.Value, 1)
		case "max":
			max = cast.ToUint64(opt.Value, 2)
		case "policy":
			policy = policyType(opt.Value)
		default:
			args = append(args, opt)
		}
	}

	result := check(service, policy, args)
	if result.Type == scaleNone {
		return
	}

	replicas := *service.Spec.Mode.Replicated.Replicas
	if result.Type == scaleUp {
		if replicas < max {
			docker.ServiceScale(service.Spec.Name, service.Version.Index, replicas+1)
			log.Get("scaler").Infof("scaler > Service '%s' scaled up for: %v", service.Spec.Name, result.Reasons)
		}
	} else if result.Type == scaleDown {
		if replicas > min {
			docker.ServiceScale(service.Spec.Name, service.Version.Index, replicas-1)
			log.Get("scaler").Infof("scaler > Service '%s' scaled down for: %v", service.Spec.Name, result.Reasons)
		}
	}
}

func check(service *swarm.Service, policy policyType, args data.Options) checkResult {
	result := checkResult{
		Reasons: make(map[string]float64),
	}
	if policy == policyAny {
		for _, arg := range args {
			st, value := checkArg(service.Spec.Name, arg)
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
			st, value := checkArg(service.Spec.Name, arg)
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

func checkArg(service string, arg data.Option) (scaleType, float64) {
	items := strings.Split(arg.Value, ":")
	if len(items) != 2 {
		log.Get("scaler").Warnf("scaler > Invalid scale argument: %s=%s", arg.Name, arg.Value)
		return scaleNone, 0
	}

	c := checkers[arg.Name]
	if c == nil {
		log.Get("scaler").Warnf("scaler > Metric checker '%s' not found", arg.Name)
		return scaleNone, 0
	}

	low := cast.ToFloat64(items[0])
	high := cast.ToFloat64(items[1])
	return c(service, low, high)
}

func cpuChecker(service string, low, high float64) (scaleType, float64) {
	query := fmt.Sprintf(`avg(rate(container_cpu_user_seconds_total{container_label_com_docker_swarm_service_name="%s"}[5m]) * 100)`, service)
	values, err := biz.Metric.GetVector(query, time.Now())
	if err != nil {
		log.Get("scaler").Error("scaler > Failed to query metrics: ", err)
		return scaleNone, 0
	}
	if len(values) == 0 {
		return scaleNone, 0
	}

	value := values[0]
	if value <= low {
		return scaleDown, value
	} else if value >= high {
		return scaleUp, value
	}
	return scaleNone, 0
}

type checkResult struct {
	Type    scaleType
	Reasons map[string]float64
}
