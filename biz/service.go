package biz

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cuigh/auxo/byte/size"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/docker"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
)

const (
	ServiceModeReplicated    = "replicated"
	ServiceModeGlobal        = "global"
	ServiceModeReplicatedJob = "replicated-job"
	ServiceModeGlobalJob     = "global-job"
)

type ServiceBiz interface {
	Search(ctx context.Context, name, mode string, pageIndex, pageSize int) (services []*ServiceBase, total int, err error)
	Find(ctx context.Context, name string, status bool) (service *Service, raw string, err error)
	Delete(ctx context.Context, name string, user web.User) (err error)
	Rollback(ctx context.Context, name string, user web.User) (err error)
	Restart(ctx context.Context, name string, user web.User) (err error)
	Scale(ctx context.Context, name string, count, version uint64, user web.User) (err error)
	Create(ctx context.Context, s *Service, user web.User) (err error)
	Update(ctx context.Context, s *Service, user web.User) (err error)
	FetchLogs(ctx context.Context, name string, lines int, timestamps bool) (stdout, stderr string, err error)
}

func NewService(d *docker.Docker, rb RegistryBiz, eb EventBiz) ServiceBiz {
	return &serviceBiz{d: d, rb: rb, eb: eb}
}

type serviceBiz struct {
	d  *docker.Docker
	rb RegistryBiz
	eb EventBiz
}

func (b *serviceBiz) Find(ctx context.Context, name string, status bool) (service *Service, raw string, err error) {
	var (
		s swarm.Service
		r []byte
	)

	s, r, err = b.d.ServiceInspect(ctx, name, status)
	if err != nil {
		if docker.IsErrNotFound(err) {
			err = nil
		}
		return
	}

	if err == nil {
		raw, err = indentJSON(r)
	}
	if err == nil {
		service = newService(&s)
		err = b.fillNetworks(ctx, service)
	}
	return
}

func (b *serviceBiz) fillNetworks(ctx context.Context, service *Service) error {
	if len(service.Endpoint.VIPs) == 0 {
		return nil
	}

	var ids = make([]string, len(service.Endpoint.VIPs))
	for i, vip := range service.Endpoint.VIPs {
		ids[i] = vip.ID
	}

	names, err := b.d.NetworkNames(ctx, ids...)
	if err == nil {
		for i := range service.Endpoint.VIPs {
			vip := &service.Endpoint.VIPs[i]
			vip.Name = names[vip.ID]
			// ingress network cannot be explicitly attached.
			if vip.Name != "ingress" {
				service.Networks = append(service.Networks, vip.Name)
			}
		}
	}
	return err
}

func (b *serviceBiz) Search(ctx context.Context, name, mode string, pageIndex, pageSize int) (services []*ServiceBase, total int, err error) {
	var list []swarm.Service
	list, total, err = b.d.ServiceList(ctx, name, mode, pageIndex, pageSize)
	if err != nil {
		return
	}

	services = make([]*ServiceBase, len(list))
	for i, s := range list {
		services[i] = newServiceBase(&s)
	}
	return
}

func (b *serviceBiz) Delete(ctx context.Context, name string, user web.User) (err error) {
	err = b.d.ServiceRemove(ctx, name)
	if err == nil {
		b.eb.CreateService(EventActionDelete, name, user)
	}
	return
}

func (b *serviceBiz) Rollback(ctx context.Context, name string, user web.User) (err error) {
	err = b.d.ServiceRollback(ctx, name)
	if err == nil {
		b.eb.CreateService(EventActionRollback, name, user)
	}
	return
}

func (b *serviceBiz) Restart(ctx context.Context, name string, user web.User) (err error) {
	err = b.d.ServiceRestart(ctx, name)
	if err == nil {
		b.eb.CreateService(EventActionRestart, name, user)
	}
	return
}

func (b *serviceBiz) Scale(ctx context.Context, name string, count, version uint64, user web.User) (err error) {
	err = b.d.ServiceScale(ctx, name, count, version)
	if err == nil {
		b.eb.CreateService(EventActionScale, name, user)
	}
	return
}

func (b *serviceBiz) Create(ctx context.Context, s *Service, user web.User) (err error) {
	spec := &swarm.ServiceSpec{TaskTemplate: swarm.TaskSpec{ContainerSpec: &swarm.ContainerSpec{}}}
	err = s.MergeTo(spec)
	if err != nil {
		return
	}

	if s.Mode == "replicated" {
		spec.Mode.Replicated = &swarm.ReplicatedService{Replicas: &s.Replicas}
	} else if s.Mode == "replicated-job" {
		spec.Mode.ReplicatedJob = &swarm.ReplicatedJob{TotalCompletions: &s.Replicas}
	} else if s.Mode == "global" {
		spec.Mode.Global = &swarm.GlobalService{}
	} else if s.Mode == "global-job" {
		spec.Mode.GlobalJob = &swarm.GlobalJob{}
	}

	auth := ""
	if i := strings.Index(s.Image, "/"); i > 0 {
		if host := s.Image[:i]; strings.Contains(host, ".") {
			auth, err = b.rb.GetAuth(ctx, host)
			if err != nil {
				return err
			}
		}
	}

	if err = b.d.ServiceCreate(ctx, spec, auth); err == nil {
		b.eb.CreateService(EventActionCreate, s.Name, user)
	}
	return
}

func (b *serviceBiz) Update(ctx context.Context, s *Service, user web.User) (err error) {
	service, _, err := b.d.ServiceInspect(ctx, s.Name, false)
	if err != nil {
		return err
	}

	spec := &service.Spec
	err = s.MergeTo(spec)
	if err != nil {
		return
	}

	if s.Mode == "replicated" && spec.Mode.Replicated != nil {
		spec.Mode.Replicated.Replicas = &s.Replicas
	} else if s.Mode == "replicated-job" && spec.Mode.ReplicatedJob != nil {
		spec.Mode.ReplicatedJob.TotalCompletions = &s.Replicas
	}

	if err = b.d.ServiceUpdate(ctx, spec, s.Version); err == nil {
		b.eb.CreateService(EventActionUpdate, s.Name, user)
	}
	return
}

func (b *serviceBiz) FetchLogs(ctx context.Context, name string, lines int, timestamps bool) (string, string, error) {
	stdout, stderr, err := b.d.ServiceLogs(ctx, name, lines, timestamps)
	if err != nil {
		return "", "", err
	}
	return stdout.String(), stderr.String(), nil
}

type ServiceBase struct {
	Name           string `json:"name"`
	Image          string `json:"image"`
	Mode           string `json:"mode"`
	Replicas       uint64 `json:"replicas"`
	DesiredTasks   uint64 `json:"desiredTasks"`
	RunningTasks   uint64 `json:"runningTasks"`
	CompletedTasks uint64 `json:"completedTasks"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type Service struct {
	ServiceBase
	ID              string       `json:"id"`
	Version         uint64       `json:"version"`
	Command         string       `json:"command"`
	Args            string       `json:"args"`
	Dir             string       `json:"dir"`
	User            string       `json:"user"`
	Hostname        string       `json:"hostname"`
	Env             data.Options `json:"env,omitempty"`
	Labels          data.Options `json:"labels,omitempty"`
	ContainerLabels data.Options `json:"containerLabels,omitempty"`
	Networks        []string     `json:"networks,omitempty"` // only for edit
	Mounts          []Mount      `json:"mounts,omitempty"`
	Update          struct {
		State   string `json:"state,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"update,omitempty"`
	Endpoint struct {
		Mode  swarm.ResolutionMode `json:"mode,omitempty"`
		Ports []EndpointPort       `json:"ports,omitempty"`
		VIPs  []EndpointVIP        `json:"vips,omitempty"`
	} `json:"endpoint"`
	Configs        []*ServiceFile `json:"configs"`
	Secrets        []*ServiceFile `json:"secrets"`
	UpdatePolicy   UpdatePolicy   `json:"updatePolicy"`
	RollbackPolicy UpdatePolicy   `json:"rollbackPolicy"`
	RestartPolicy  RestartPolicy  `json:"restartPolicy"`
	Resource       struct {
		Limit   ServiceResource `json:"limit"`
		Reserve ServiceResource `json:"reserve"`
	} `json:"resource"`
	Placement struct {
		Constraints []PlacementConstraint `json:"constraints"`
		Preferences []string              `json:"preferences"`
		//Platforms   []Platform            `json:"platforms"`
	} `json:"placement"`
	LogDriver struct {
		Name    string       `json:"name"`
		Options data.Options `json:"options"`
	} `json:"logDriver"`
	DNS struct {
		Servers []string `json:"servers"`
		Search  []string `json:"search"`
		Options []string `json:"options"`
	} `json:"dns"`
	Hosts []string `json:"hosts"`
}

type PlacementConstraint struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Operator string `json:"op"`
}

func NewPlacementConstraint(c string) PlacementConstraint {
	var (
		pc    = PlacementConstraint{}
		items []string
	)

	if items = strings.SplitN(c, "==", 2); len(items) == 2 {
		pc.Operator = "=="
	} else if items = strings.SplitN(c, "!=", 2); len(items) == 2 {
		pc.Operator = "!="
	}
	if pc.Operator != "" {
		pc.Name = strings.TrimSpace(items[0])
		pc.Value = strings.TrimSpace(items[1])
	}

	return pc
}

func (pc *PlacementConstraint) ToConstraint() string {
	if pc.Name != "" && pc.Value != "" && pc.Operator != "" {
		return fmt.Sprintf("%s %s %s", pc.Name, pc.Operator, pc.Value)
	}
	return ""
}

type ServiceResource struct {
	CPU    float64 `json:"cpu,omitempty"`
	Memory string  `json:"memory,omitempty"`
}

func newServiceResourceFromResources(res *swarm.Resources) ServiceResource {
	ri := ServiceResource{}
	if res != nil {
		ri.CPU = float64(res.NanoCPUs) / 1e9
		if res.MemoryBytes > 0 {
			ri.Memory = size.Size(res.MemoryBytes).String()
		}
	}
	return ri
}

func newServiceResourceFromLimit(res *swarm.Limit) ServiceResource {
	ri := ServiceResource{}
	if res != nil {
		ri.CPU = float64(res.NanoCPUs) / 1e9
		if res.MemoryBytes > 0 {
			ri.Memory = size.Size(res.MemoryBytes).String()
		}
	}
	return ri
}

func (r ServiceResource) IsSet() bool {
	return r.CPU > 0 || r.Memory != ""
}

func (r ServiceResource) ToResources() (res *swarm.Resources, err error) {
	res = &swarm.Resources{
		NanoCPUs: int64(r.CPU * 1e9),
	}
	if r.Memory != "" {
		var s size.Size
		if s, err = size.Parse(r.Memory); err != nil {
			return nil, err
		}
		res.MemoryBytes = int64(s)
	}
	return
}

func (r ServiceResource) ToLimit() (res *swarm.Limit, err error) {
	res = &swarm.Limit{
		NanoCPUs: int64(r.CPU * 1e9),
	}
	if r.Memory != "" {
		var s size.Size
		if s, err = size.Parse(r.Memory); err != nil {
			return nil, err
		}
		res.MemoryBytes = int64(s)
	}
	return
}

// EndpointPort represents the config of a port.
type EndpointPort struct {
	Name          string                      `json:"name,omitempty"`
	Protocol      swarm.PortConfigProtocol    `json:"protocol,omitempty"`
	TargetPort    uint32                      `json:"targetPort,omitempty"`
	PublishedPort uint32                      `json:"publishedPort,omitempty"`
	PublishMode   swarm.PortConfigPublishMode `json:"mode,omitempty"`
}

// EndpointVIP represents the virtual ip of a port.
type EndpointVIP struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	IP   string `json:"ip,omitempty"`
}

type ServiceFile struct {
	//ID   string      `json:"id"`
	//Name string      `json:"name"`
	Key  string `json:"key"` // ID:Name
	Path string `json:"path"`
	GID  string `json:"gid"`
	UID  string `json:"uid"`
	Mode uint32 `json:"mode"`
}

func newServiceFile(id, name string, filename, uid, gid string, mode os.FileMode) *ServiceFile {
	m, _ := strconv.ParseUint(strconv.FormatUint(uint64(mode), 8), 10, 32)
	return &ServiceFile{
		Key:  id + ":" + name,
		Path: filename,
		UID:  uid,
		GID:  gid,
		Mode: uint32(m),
	}
}

func (f *ServiceFile) ToConfig() *swarm.ConfigReference {
	mode, _ := strconv.ParseUint(strconv.FormatUint(uint64(f.Mode), 10), 8, 32)
	pair := strings.Split(f.Key, ":")
	return &swarm.ConfigReference{
		ConfigID:   pair[0],
		ConfigName: pair[1],
		File: &swarm.ConfigReferenceFileTarget{
			Name: f.Path,
			UID:  f.UID,
			GID:  f.GID,
			Mode: os.FileMode(mode),
		},
	}
}

func (f *ServiceFile) ToSecret() *swarm.SecretReference {
	mode, _ := strconv.ParseUint(strconv.FormatUint(uint64(f.Mode), 10), 8, 32)
	pair := strings.Split(f.Key, ":")
	return &swarm.SecretReference{
		SecretID:   pair[0],
		SecretName: pair[1],
		File: &swarm.SecretReferenceFileTarget{
			Name: f.Path,
			UID:  f.UID,
			GID:  f.GID,
			Mode: os.FileMode(mode),
		},
	}
}

// Mount represents a mount (volume).
type Mount struct {
	Type        mount.Type        `json:"type,omitempty"`
	Source      string            `json:"source,omitempty"`
	Target      string            `json:"target,omitempty"`
	Readonly    bool              `json:"readonly,omitempty"`
	Consistency mount.Consistency `json:"consistency,omitempty"`
	//BindOptions   *BindOptions   `json:",omitempty"`
	//VolumeOptions *VolumeOptions `json:",omitempty"`
	//TmpfsOptions  *TmpfsOptions  `json:",omitempty"`
}

type UpdatePolicy struct {
	Parallelism   uint64 `json:"parallelism,omitempty"`
	Delay         string `json:"delay,omitempty"`
	FailureAction string `json:"failureAction,omitempty"`
	Order         string `json:"order,omitempty"`
	//Monitor time.Duration `json:",omitempty"`
	//MaxFailureRatio float32
}

func (p *UpdatePolicy) Convert() *swarm.UpdateConfig {
	if p.Parallelism == 0 && p.Delay == "" && p.FailureAction == "" && p.Order == "" {
		return nil
	}

	delay, _ := time.ParseDuration(p.Delay)
	return &swarm.UpdateConfig{
		Parallelism:   p.Parallelism,
		Delay:         delay,
		FailureAction: p.FailureAction,
		Order:         p.Order,
	}
}

func newUpdatePolicy(c *swarm.UpdateConfig) UpdatePolicy {
	p := UpdatePolicy{}
	if c != nil {
		p.Parallelism = c.Parallelism
		p.Delay = c.Delay.String()
		p.FailureAction = c.FailureAction
		p.Order = c.Order
	}
	return p
}

type RestartPolicy struct {
	Condition   swarm.RestartPolicyCondition `json:"condition,omitempty"`
	Delay       string                       `json:"delay,omitempty"`
	MaxAttempts uint64                       `json:"maxAttempts,omitempty"`
	Window      string                       `json:"window,omitempty"`
}

func (p *RestartPolicy) Convert() *swarm.RestartPolicy {
	if p.MaxAttempts == 0 && p.Delay == "" && p.Condition == "" && p.Window == "" {
		return nil
	}

	policy := &swarm.RestartPolicy{Condition: p.Condition}
	if delay, err := time.ParseDuration(p.Delay); err == nil {
		policy.Delay = &delay
	}
	if window, err := time.ParseDuration(p.Window); err == nil {
		policy.Window = &window
	}
	if p.MaxAttempts > 0 {
		policy.MaxAttempts = &p.MaxAttempts
	}
	return policy
}

func newRestartPolicy(p *swarm.RestartPolicy) RestartPolicy {
	policy := RestartPolicy{}
	if p != nil {
		policy.Condition = p.Condition
		if p.Delay != nil {
			policy.Delay = p.Delay.String()
		}
		if p.MaxAttempts != nil {
			policy.MaxAttempts = *p.MaxAttempts
		}
		if p.Window != nil {
			policy.Window = p.Window.String()
		}
	}
	return policy
}

func newServiceBase(s *swarm.Service) *ServiceBase {
	service := &ServiceBase{
		Name:      s.Spec.Name,
		Image:     normalizeImage(s.Spec.TaskTemplate.ContainerSpec.Image),
		CreatedAt: formatTime(s.CreatedAt),
		UpdatedAt: formatTime(s.UpdatedAt),
	}

	if s.ServiceStatus != nil {
		service.RunningTasks = s.ServiceStatus.RunningTasks
		service.DesiredTasks = s.ServiceStatus.DesiredTasks
		service.CompletedTasks = s.ServiceStatus.CompletedTasks
	}

	if s.Spec.Mode.Replicated != nil {
		service.Mode = ServiceModeReplicated
		service.Replicas = *s.Spec.Mode.Replicated.Replicas
	} else if s.Spec.Mode.Global != nil {
		service.Mode = ServiceModeGlobal
	} else if s.Spec.Mode.ReplicatedJob != nil {
		service.Mode = ServiceModeReplicatedJob
		service.Replicas = *s.Spec.Mode.ReplicatedJob.TotalCompletions
	} else if s.Spec.Mode.GlobalJob != nil {
		service.Mode = ServiceModeGlobalJob
	}
	return service
}

func newService(s *swarm.Service) *Service {
	service := &Service{
		ServiceBase:     *newServiceBase(s),
		ID:              s.ID,
		Version:         s.Version.Index,
		Env:             envToOptions(s.Spec.TaskTemplate.ContainerSpec.Env),
		Labels:          mapToOptions(s.Spec.Labels),
		ContainerLabels: mapToOptions(s.Spec.TaskTemplate.ContainerSpec.Labels),
		Command:         strings.Join(s.Spec.TaskTemplate.ContainerSpec.Command, " "),
		Args:            strings.Join(s.Spec.TaskTemplate.ContainerSpec.Args, " "),
		Dir:             s.Spec.TaskTemplate.ContainerSpec.Dir,
		User:            s.Spec.TaskTemplate.ContainerSpec.User,
		UpdatePolicy:    newUpdatePolicy(s.Spec.UpdateConfig),
		RollbackPolicy:  newUpdatePolicy(s.Spec.RollbackConfig),
		RestartPolicy:   newRestartPolicy(s.Spec.TaskTemplate.RestartPolicy),
	}

	if s.UpdateStatus != nil {
		service.Update.State = string(s.UpdateStatus.State)
		service.Update.Message = s.UpdateStatus.Message
	}

	// Endpoint
	service.Endpoint.Mode = s.Endpoint.Spec.Mode
	for _, vip := range s.Endpoint.VirtualIPs {
		service.Endpoint.VIPs = append(service.Endpoint.VIPs, EndpointVIP{
			ID: vip.NetworkID,
			IP: vip.Addr,
		})
	}
	for _, p := range s.Endpoint.Ports {
		service.Endpoint.Ports = append(service.Endpoint.Ports, EndpointPort{
			Name:          p.Name,
			Protocol:      p.Protocol,
			TargetPort:    p.TargetPort,
			PublishedPort: p.PublishedPort,
			PublishMode:   p.PublishMode,
		})
	}

	// Mounts
	if l := len(s.Spec.TaskTemplate.ContainerSpec.Mounts); l > 0 {
		service.Mounts = make([]Mount, l)
		for i, m := range s.Spec.TaskTemplate.ContainerSpec.Mounts {
			service.Mounts[i] = Mount{
				Type:        m.Type,
				Source:      m.Source,
				Target:      m.Target,
				Readonly:    m.ReadOnly,
				Consistency: m.Consistency,
			}
		}
	}

	// Files
	for _, f := range s.Spec.TaskTemplate.ContainerSpec.Configs {
		file := newServiceFile(f.ConfigID, f.ConfigName, f.File.Name, f.File.UID, f.File.GID, f.File.Mode)
		service.Configs = append(service.Configs, file)
	}
	for _, f := range s.Spec.TaskTemplate.ContainerSpec.Secrets {
		file := newServiceFile(f.SecretID, f.SecretName, f.File.Name, f.File.UID, f.File.GID, f.File.Mode)
		service.Secrets = append(service.Secrets, file)
	}

	// Resource
	if s.Spec.TaskTemplate.Resources != nil {
		service.Resource.Limit = newServiceResourceFromLimit(s.Spec.TaskTemplate.Resources.Limits)
		service.Resource.Reserve = newServiceResourceFromResources(s.Spec.TaskTemplate.Resources.Reservations)
	}

	// Placement
	if s.Spec.TaskTemplate.Placement != nil {
		for _, c := range s.Spec.TaskTemplate.Placement.Constraints {
			service.Placement.Constraints = append(service.Placement.Constraints, NewPlacementConstraint(c))
		}
		for _, p := range s.Spec.TaskTemplate.Placement.Preferences {
			service.Placement.Preferences = append(service.Placement.Preferences, p.Spread.SpreadDescriptor)
		}
	}

	// LogDriver
	if s.Spec.TaskTemplate.LogDriver != nil {
		service.LogDriver.Name = s.Spec.TaskTemplate.LogDriver.Name
		service.LogDriver.Options = mapToOptions(s.Spec.TaskTemplate.LogDriver.Options)
	}

	// DNS
	if s.Spec.TaskTemplate.ContainerSpec.DNSConfig != nil {
		service.DNS.Servers = s.Spec.TaskTemplate.ContainerSpec.DNSConfig.Nameservers
		service.DNS.Search = s.Spec.TaskTemplate.ContainerSpec.DNSConfig.Search
		service.DNS.Options = s.Spec.TaskTemplate.ContainerSpec.DNSConfig.Options
	}
	return service
}

func (s *Service) MergeTo(spec *swarm.ServiceSpec) (err error) {
	spec.Name = s.Name
	spec.Labels = toMap(s.Labels)
	spec.TaskTemplate.ContainerSpec.Image = s.Image
	spec.TaskTemplate.ContainerSpec.Dir = s.Dir
	spec.TaskTemplate.ContainerSpec.User = s.User
	spec.TaskTemplate.ContainerSpec.Hostname = s.Hostname
	spec.TaskTemplate.ContainerSpec.Hosts = s.Hosts
	spec.TaskTemplate.ContainerSpec.Labels = toMap(s.ContainerLabels)
	spec.TaskTemplate.ContainerSpec.Env = toEnv(s.Env)
	spec.TaskTemplate.ContainerSpec.Command = parseArgs(s.Command)
	spec.TaskTemplate.ContainerSpec.Args = parseArgs(s.Args)

	// Networks
	spec.TaskTemplate.Networks = nil
	for _, n := range s.Networks {
		spec.TaskTemplate.Networks = append(spec.TaskTemplate.Networks, swarm.NetworkAttachmentConfig{Target: n})
	}

	// Endpoint
	if s.Endpoint.Mode == "" && len(s.Endpoint.Ports) == 0 {
		spec.EndpointSpec = nil
	} else {
		spec.EndpointSpec = &swarm.EndpointSpec{Mode: s.Endpoint.Mode}
		for _, p := range s.Endpoint.Ports {
			spec.EndpointSpec.Ports = append(spec.EndpointSpec.Ports, swarm.PortConfig{
				Name:          p.Name,
				Protocol:      p.Protocol,
				TargetPort:    p.TargetPort,
				PublishedPort: p.PublishedPort,
				PublishMode:   p.PublishMode,
			})
		}
	}

	// Mounts
	spec.TaskTemplate.ContainerSpec.Mounts = nil
	for _, m := range s.Mounts {
		spec.TaskTemplate.ContainerSpec.Mounts = append(spec.TaskTemplate.ContainerSpec.Mounts, mount.Mount{
			Type:        m.Type,
			Source:      m.Source,
			Target:      m.Target,
			ReadOnly:    m.Readonly,
			Consistency: m.Consistency,
		})
	}

	// Configs
	spec.TaskTemplate.ContainerSpec.Configs = nil
	for _, f := range s.Configs {
		spec.TaskTemplate.ContainerSpec.Configs = append(spec.TaskTemplate.ContainerSpec.Configs, f.ToConfig())
	}

	// Secrets
	spec.TaskTemplate.ContainerSpec.Secrets = nil
	for _, f := range s.Secrets {
		spec.TaskTemplate.ContainerSpec.Secrets = append(spec.TaskTemplate.ContainerSpec.Secrets, f.ToSecret())
	}

	// Resource
	if s.Resource.Limit.IsSet() || s.Resource.Reserve.IsSet() {
		spec.TaskTemplate.Resources = &swarm.ResourceRequirements{}
		if s.Resource.Limit.IsSet() {
			spec.TaskTemplate.Resources.Limits, err = s.Resource.Limit.ToLimit()
			if err != nil {
				return
			}
		}
		if s.Resource.Limit.IsSet() {
			spec.TaskTemplate.Resources.Reservations, err = s.Resource.Reserve.ToResources()
			if err != nil {
				return
			}
		}
	} else {
		spec.TaskTemplate.Resources = nil
	}

	// Placement
	if len(s.Placement.Constraints) == 0 && len(s.Placement.Preferences) == 0 {
		spec.TaskTemplate.Placement = nil
	} else {
		spec.TaskTemplate.Placement = &swarm.Placement{}
		for _, c := range s.Placement.Constraints {
			if cons := c.ToConstraint(); cons != "" {
				spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, cons)
			}
		}
		for _, p := range s.Placement.Preferences {
			if p != "" {
				pref := swarm.PlacementPreference{Spread: &swarm.SpreadOver{SpreadDescriptor: p}}
				spec.TaskTemplate.Placement.Preferences = append(spec.TaskTemplate.Placement.Preferences, pref)
			}
		}
	}

	// Schedule
	spec.UpdateConfig = s.UpdatePolicy.Convert()
	spec.RollbackConfig = s.RollbackPolicy.Convert()
	spec.TaskTemplate.RestartPolicy = s.RestartPolicy.Convert()

	// LogDriver
	if s.LogDriver.Name == "" {
		spec.TaskTemplate.LogDriver = nil
	} else {
		spec.TaskTemplate.LogDriver = &swarm.Driver{
			Name:    s.LogDriver.Name,
			Options: toMap(s.LogDriver.Options),
		}
	}

	// Host & DNS
	if len(s.DNS.Options) == 0 && len(s.DNS.Options) == 0 && len(s.DNS.Options) == 0 {
		spec.TaskTemplate.ContainerSpec.DNSConfig = nil
	} else {
		spec.TaskTemplate.ContainerSpec.DNSConfig = &swarm.DNSConfig{
			Nameservers: s.DNS.Servers,
			Search:      s.DNS.Search,
			Options:     s.DNS.Options,
		}
	}

	return nil
}
