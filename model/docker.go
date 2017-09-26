package model

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
	"fmt"
	"strconv"
	"os"

	"github.com/cuigh/auxo/data/size"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
)

type Registry struct {
	ID        string    `bson:"_id" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name,omitempty"`
	URL       string    `bson:"url" json:"url,omitempty"`
	Username  string    `bson:"username" json:"username,omitempty"`
	Password  string    `bson:"password" json:"password,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}

func (r *Registry) Match(image string) bool {
	return strings.HasPrefix(image, r.URL)
}

func (r *Registry) GetEncodedAuth() string {
	cfg := &types.AuthConfig{
		ServerAddress: r.URL,
		Username:      r.Username,
		Password:      r.Password,
	}
	if buf, e := json.Marshal(cfg); e == nil {
		return base64.URLEncoding.EncodeToString(buf)
	}
	return ""
}

type Option struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Options []*Option

func NewOptions(m map[string]string) Options {
	if len(m) == 0 {
		return nil
	}

	opts := Options{}
	for k, v := range m {
		opts = append(opts, &Option{Name: k, Value: v})
	}
	return opts
}

func (opts Options) ToMap() map[string]string {
	if len(opts) == 0 {
		return nil
	}

	m := make(map[string]string)
	for _, opt := range opts {
		if opt != nil && opt.Name != "" && opt.Value != "" {
			m[opt.Name] = opt.Value
		}
	}
	return m
}

type StackListInfo struct {
	Name     string
	Services []string
}

type ServiceListInfo struct {
	Name      string
	Image     string
	Mode      string
	Actives   uint64
	Replicas  uint64
	UpdatedAt time.Time
}

func NewServiceListInfo(service swarm.Service, actives uint64) *ServiceListInfo {
	info := &ServiceListInfo{
		Name:      service.Spec.Name,
		Image:     normalizeImage(service.Spec.TaskTemplate.ContainerSpec.Image),
		Actives:   actives,
		UpdatedAt: service.UpdatedAt.Local(),
	}
	if service.Spec.Mode.Replicated != nil && service.Spec.Mode.Replicated.Replicas != nil {
		info.Mode = "replicated"
		info.Replicas = *service.Spec.Mode.Replicated.Replicas
	} else if service.Spec.Mode.Global != nil {
		info.Mode = "global"
		// info.Replicas = tasksNoShutdown[service.ID]
	}
	return info
}

type ServiceDetailInfo struct {
	swarm.Service
	// Name      string
	// Image     string
	// Mode      string
	// Actives   uint64
	Replicas uint64
	// UpdatedAt time.Time
	Env      map[string]string
	Networks []Network
	Command  string
	Args     string
}

type Network struct {
	ID      string
	Name    string
	Address string
}

func NewServiceDetailInfo(service swarm.Service) *ServiceDetailInfo {
	info := &ServiceDetailInfo{
		Service: service,
		Command: strings.Join(service.Spec.TaskTemplate.ContainerSpec.Command, " "),
		Args:    strings.Join(service.Spec.TaskTemplate.ContainerSpec.Args, " "),
	}
	if service.Spec.Mode.Replicated != nil {
		info.Replicas = *service.Spec.Mode.Replicated.Replicas
	}
	if len(service.Spec.TaskTemplate.ContainerSpec.Env) > 0 {
		info.Env = make(map[string]string)
		for _, env := range service.Spec.TaskTemplate.ContainerSpec.Env {
			pair := strings.SplitN(env, "=", 2)
			if len(pair) == 2 {
				info.Env[pair[0]] = pair[1]
			}
		}
	}
	return info
}

type ServiceInfo struct {
	Name            string       `json:"name"`
	Registry        string       `json:"registry"`
	RegistryAuth    string       `json:"-"`
	Image           string       `json:"image"`
	Mode            string       `json:"mode"`
	Replicas        uint64       `json:"replicas"`
	Command         string       `json:"command"`
	Args            string       `json:"args"`
	Dir             string       `json:"dir"`
	User            string       `json:"user"`
	Networks        []string     `json:"networks"`
	Secrets         []ConfigInfo `json:"secrets"`
	Configs         []ConfigInfo `json:"configs"`
	Environments    Options      `json:"envs"`
	ServiceLabels   Options      `json:"slabels"`
	ContainerLabels Options      `json:"clabels"`
	Endpoint        struct {
		Mode  swarm.ResolutionMode `json:"mode"`
		Ports []EndpointPort       `json:"ports"`
	} `json:"endpoint"`
	Mounts    []Mount `json:"mounts"`
	LogDriver struct {
		Name    string  `json:"name"`
		Options Options `json:"options"`
	} `json:"log_driver"`
	Placement struct {
		Constraints []PlacementConstraint `json:"constraints"`
		Preferences []PlacementPreference `json:"preferences"`
	} `json:"placement"`
	UpdatePolicy struct {
		// Maximum number of tasks updated simultaneously (0 to update all at once), default 1.
		Parallelism uint64 `json:"parallelism"`
		// Amount of time between updates.
		Delay string `json:"delay,omitempty"`
		// FailureAction is the action to take when an update failures. (“pause”|“continue”|“rollback”) (default “pause”)
		FailureAction string `json:"failure_action"`
		// Update order (“start-first”|“stop-first”) (default “stop-first”)
		Order string `json:"order"`
	} `json:"update_policy"`
	RollbackPolicy struct {
		// Maximum number of tasks updated simultaneously (0 to update all at once), default 1.
		Parallelism uint64 `json:"parallelism"`
		// Amount of time between updates.
		Delay string `json:"delay,omitempty"`
		// FailureAction is the action to take when an update failures. (“pause”|“continue”) (default “pause”)
		FailureAction string `json:"failure_action"`
		// Update order (“start-first”|“stop-first”) (default “stop-first”)
		Order string `json:"order"`
	} `json:"rollback_policy"`
	RestartPolicy struct {
		// Restart when condition is met (“none”|“on-failure”|“any”) (default “any”)
		Condition swarm.RestartPolicyCondition `json:"condition"`
		// Maximum number of restarts before giving up
		MaxAttempts uint64 `json:"max_attempts"`
		// Delay between restart attempts
		Delay string `json:"delay,omitempty"`
		// Window used to evaluate the restart policy.
		Window string `json:"window,omitempty"`
	} `json:"restart_policy"`
	Resource struct {
		Limit   ResourceInfo `json:"limit"`
		Reserve ResourceInfo `json:"reserve"`
	} `json:"resource"`
}

func NewServiceInfo(service swarm.Service) *ServiceInfo {
	spec := service.Spec
	si := &ServiceInfo{
		Name: spec.Name,
		//Hostname:    spec.TaskTemplate.ContainerSpec.Hostname,
		Image:           spec.TaskTemplate.ContainerSpec.Image,
		Command:         strings.Join(spec.TaskTemplate.ContainerSpec.Command, " "),
		Args:            strings.Join(spec.TaskTemplate.ContainerSpec.Args, " "),
		Dir:             spec.TaskTemplate.ContainerSpec.Dir,
		User:            spec.TaskTemplate.ContainerSpec.User,
		ServiceLabels:   NewOptions(spec.Labels),
		ContainerLabels: NewOptions(spec.TaskTemplate.ContainerSpec.Labels),
	}
	for _, vip := range service.Endpoint.VirtualIPs {
		si.Networks = append(si.Networks, vip.NetworkID)
	}
	if spec.EndpointSpec != nil {
		si.Endpoint.Mode = spec.EndpointSpec.Mode
		for _, p := range spec.EndpointSpec.Ports {
			si.Endpoint.Ports = append(si.Endpoint.Ports, EndpointPort{
				Protocol:      p.Protocol,
				TargetPort:    p.TargetPort,
				PublishedPort: p.PublishedPort,
				PublishMode:   p.PublishMode,
			})
		}
	}
	if spec.Mode.Global != nil {
		si.Mode = "global"
	} else if spec.Mode.Replicated != nil {
		si.Mode, si.Replicas = "replicated", *spec.Mode.Replicated.Replicas
	}
	if len(spec.TaskTemplate.ContainerSpec.Env) > 0 {
		si.Environments = Options{}
		for _, env := range spec.TaskTemplate.ContainerSpec.Env {
			pair := strings.SplitN(env, "=", 2)
			si.Environments = append(si.Environments, &Option{Name: pair[0], Value: pair[1]})
		}
	}
	for _, m := range spec.TaskTemplate.ContainerSpec.Mounts {
		mnt := Mount{
			Type:     m.Type,
			Source:   m.Source,
			Target:   m.Target,
			ReadOnly: m.ReadOnly,
		}
		if m.BindOptions != nil {
			mnt.Propagation = m.BindOptions.Propagation
		}
		si.Mounts = append(si.Mounts, mnt)
	}
	if spec.TaskTemplate.Resources != nil {
		si.Resource.Limit = NewResourceInfo(spec.TaskTemplate.Resources.Limits)
		si.Resource.Reserve = NewResourceInfo(spec.TaskTemplate.Resources.Reservations)
	}
	if spec.TaskTemplate.LogDriver != nil {
		si.LogDriver.Name = spec.TaskTemplate.LogDriver.Name
		si.LogDriver.Options = NewOptions(spec.TaskTemplate.LogDriver.Options)
	}
	if spec.TaskTemplate.Placement != nil {
		for _, c := range spec.TaskTemplate.Placement.Constraints {
			si.Placement.Constraints = append(si.Placement.Constraints, NewPlacementConstraint(c))
		}
		for _, p := range spec.TaskTemplate.Placement.Preferences {
			si.Placement.Preferences = append(si.Placement.Preferences, PlacementPreference{Spread: p.Spread.SpreadDescriptor})
		}
	}
	if spec.UpdateConfig != nil {
		si.UpdatePolicy.Parallelism = spec.UpdateConfig.Parallelism
		si.UpdatePolicy.Delay = spec.UpdateConfig.Delay.String()
		si.UpdatePolicy.FailureAction = spec.UpdateConfig.FailureAction
		si.UpdatePolicy.Order = spec.UpdateConfig.Order
	}
	if spec.RollbackConfig != nil {
		si.RollbackPolicy.Parallelism = spec.RollbackConfig.Parallelism
		si.RollbackPolicy.Delay = spec.RollbackConfig.Delay.String()
		si.RollbackPolicy.FailureAction = spec.RollbackConfig.FailureAction
		si.RollbackPolicy.Order = spec.RollbackConfig.Order
	}
	if spec.TaskTemplate.RestartPolicy != nil {
		si.RestartPolicy.Condition = spec.TaskTemplate.RestartPolicy.Condition
		if spec.TaskTemplate.RestartPolicy.MaxAttempts != nil {
			si.RestartPolicy.MaxAttempts = *spec.TaskTemplate.RestartPolicy.MaxAttempts
		}
		if spec.TaskTemplate.RestartPolicy.Delay != nil {
			si.RestartPolicy.Delay = spec.TaskTemplate.RestartPolicy.Delay.String()
		}
		if spec.TaskTemplate.RestartPolicy.Window != nil {
			si.RestartPolicy.Window = spec.TaskTemplate.RestartPolicy.Window.String()
		}
	}
	for _, s := range spec.TaskTemplate.ContainerSpec.Secrets {
		secret := NewSecretInfo(s)
		si.Secrets = append(si.Secrets, secret)
	}
	for _, c := range spec.TaskTemplate.ContainerSpec.Configs {
		config := NewConfigInfo(c)
		si.Configs = append(si.Configs, config)
	}
	return si
}

func (si *ServiceInfo) ToServiceSpec() swarm.ServiceSpec {
	return swarm.ServiceSpec{}
}

type EndpointPort struct {
	Name     string                   `json:"name"`
	Protocol swarm.PortConfigProtocol `json:"protocol"`
	// TargetPort is the port inside the container
	TargetPort uint32 `json:"target_port"`
	// PublishedPort is the port on the swarm hosts
	PublishedPort uint32 `json:"published_port"`
	// PublishMode is the mode in which port is published
	PublishMode swarm.PortConfigPublishMode `json:"publish_mode"`
}

type Mount struct {
	Type        mount.Type        `json:"type"`
	Source      string            `json:"src"`
	Target      string            `json:"dst"`
	ReadOnly    bool              `json:"read_only"`
	Propagation mount.Propagation `json:"propagation"`
}

type ResourceInfo struct {
	CPU    float64 `json:"cpu"`
	Memory string  `json:"memory"`
}

func NewResourceInfo(res *swarm.Resources) ResourceInfo {
	ri := ResourceInfo{}
	if res != nil {
		ri.CPU = float64(res.NanoCPUs) / 1e9
		ri.Memory = size.Size(res.MemoryBytes).String()
	}
	return ri
}

func (r ResourceInfo) IsSet() bool {
	return r.CPU > 0 || r.Memory != ""
}

func (r ResourceInfo) ToResources() (res *swarm.Resources, err error) {
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

type ConfigInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FileName string `json:"file_name"`
	UID      string `json:"uid"`
	GID      string `json:"gid"`
	Mode     uint32 `json:"mode"`
}

func NewSecretInfo(ref *swarm.SecretReference) ConfigInfo {
	mode, _ := strconv.ParseUint(strconv.FormatUint(uint64(ref.File.Mode), 8), 10, 32)
	return ConfigInfo{
		ID:       ref.SecretID,
		Name:     ref.SecretName,
		FileName: ref.File.Name,
		UID:      ref.File.UID,
		GID:      ref.File.GID,
		Mode:     uint32(mode),
	}
}

func NewConfigInfo(ref *swarm.ConfigReference) ConfigInfo {
	mode, _ := strconv.ParseUint(strconv.FormatUint(uint64(ref.File.Mode), 8), 10, 32)
	return ConfigInfo{
		ID:       ref.ConfigID,
		Name:     ref.ConfigName,
		FileName: ref.File.Name,
		UID:      ref.File.UID,
		GID:      ref.File.GID,
		Mode:     uint32(mode),
	}
}

func (ci ConfigInfo) ToSecret() *swarm.SecretReference {
	mode, _ := strconv.ParseUint(strconv.FormatUint(uint64(ci.Mode), 10), 8, 32)
	return &swarm.SecretReference{
		SecretID:   ci.ID,
		SecretName: ci.Name,
		File: &swarm.SecretReferenceFileTarget{
			Name: ci.FileName,
			UID:  ci.UID,
			GID:  ci.GID,
			Mode: os.FileMode(mode),
		},
	}
}

func (ci ConfigInfo) ToConfig() *swarm.ConfigReference {
	mode, _ := strconv.ParseUint(strconv.FormatUint(uint64(ci.Mode), 10), 8, 32)
	return &swarm.ConfigReference{
		ConfigID:   ci.ID,
		ConfigName: ci.Name,
		File: &swarm.ConfigReferenceFileTarget{
			Name: ci.FileName,
			UID:  ci.UID,
			GID:  ci.GID,
			Mode: os.FileMode(mode),
		},
	}
}

type PlacementPreference struct {
	Spread string `json:"spread"`
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

type TaskInfo struct {
	swarm.Task
	NodeName string
	Image    string
}

func NewTaskInfo(task swarm.Task, nodeName string) *TaskInfo {
	return &TaskInfo{
		Task:     task,
		NodeName: nodeName,
		Image:    normalizeImage(task.Spec.ContainerSpec.Image),
	}
}

type NodeListInfo struct {
	ID      string
	Name    string
	Role    swarm.NodeRole
	Version string
	CPU     int64
	Memory  float32
	Address string
	Status  swarm.NodeState
	Leader  bool
}

func NewNodeListInfo(node swarm.Node) *NodeListInfo {
	info := &NodeListInfo{
		ID:      node.ID,
		Name:    node.Spec.Name,
		Role:    node.Spec.Role,
		Version: node.Description.Engine.EngineVersion,
		CPU:     node.Description.Resources.NanoCPUs / 1e9,
		Memory:  float32(node.Description.Resources.MemoryBytes>>20) / 1024,
		Address: node.Status.Addr,
		Status:  node.Status.State,
		Leader:  node.ManagerStatus != nil && node.ManagerStatus.Leader,
	}
	if info.Name == "" {
		info.Name = node.Description.Hostname
	}
	return info
}

type NodeUpdateInfo struct {
	Version      uint64                 `json:"version"`
	Name         string                 `json:"name"`
	Role         swarm.NodeRole         `json:"role"`
	Availability swarm.NodeAvailability `json:"availability"`
	Labels       Options                `json:"labels"`
}

type NetworkCreateInfo struct {
	Name         string `json:"name"`
	Driver       string `json:"driver"`
	CustomDriver string `json:"custom_driver"`
	Internal     bool   `json:"internal"`
	Attachable   bool
	IPV4         struct {
		Subnet  string `json:"subnet"`
		Gateway string `json:"gateway"`
		IPRange string `json:"ip_range"`
	} `json:"ipv4"`
	IPV6 struct {
		Enabled bool   `json:"enabled"`
		Subnet  string `json:"subnet"`
		Gateway string `json:"gateway"`
	} `json:"ipv6"`
	Options   Options `json:"options"`
	Labels    Options `json:"labels"`
	LogDriver struct {
		Name    string  `json:"name"`
		Options Options `json:"options"`
	} `json:"log_driver"`
}

type VolumeCreateInfo struct {
	Name         string  `json:"name"`
	Driver       string  `json:"driver"`
	CustomDriver string  `json:"custom_driver"`
	Options      Options `json:"options"`
	Labels       Options `json:"labels"`
}

type ImageListInfo struct {
	types.ImageSummary
	CreatedAt time.Time
}

func NewImageListInfo(image types.ImageSummary) *ImageListInfo {
	info := &ImageListInfo{
		ImageSummary: image,
		CreatedAt:    time.Unix(image.Created, 0),
	}
	return info
}

type ContainerListInfo struct {
	types.Container
	CreatedAt time.Time
}

func NewContainerListInfo(container types.Container) *ContainerListInfo {
	info := &ContainerListInfo{
		Container: container,
		CreatedAt: time.Unix(container.Created, 0),
	}
	return info
}

func normalizeImage(image string) string {
	// remove hash added by docker
	if i := strings.Index(image, "@sha256:"); i > 0 {
		image = image[:i]
	}
	return image
}
