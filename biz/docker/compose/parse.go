package compose

import (
	"fmt"
	"os"
	"sort"
	"strings"

	composetypes "github.com/cuigh/swirl/biz/docker/compose/types"
)

func Parse(name, content string) (*composetypes.Config, error) {
	//absPath, err := filepath.Abs(composefile)
	//if err != nil {
	//	return details, err
	//}
	//details.WorkingDir = filepath.Dir(absPath)

	configFile, err := getConfigFile(name, content)
	if err != nil {
		return nil, err
	}

	env, err := buildEnvironment(os.Environ())
	if err != nil {
		return nil, err
	}

	details := composetypes.ConfigDetails{
		ConfigFiles: []composetypes.ConfigFile{*configFile},
		Environment: env,
	}
	cfg, err := Load(details)
	if err != nil {
		if fpe, ok := err.(*ForbiddenPropertiesError); ok {
			err = fmt.Errorf("Compose file contains unsupported options:\n\n%s\n",
				propertyWarnings(fpe.Properties))
		}
	}
	return cfg, err
}

func propertyWarnings(properties map[string]string) string {
	var msgs []string
	for name, description := range properties {
		msgs = append(msgs, fmt.Sprintf("%s: %s", name, description))
	}
	sort.Strings(msgs)
	return strings.Join(msgs, "\n\n")
}

func getConfigFile(name, content string) (*composetypes.ConfigFile, error) {
	config, err := ParseYAML([]byte(content))
	if err != nil {
		return nil, err
	}
	return &composetypes.ConfigFile{
		Filename: name,
		Config:   config,
	}, nil
}

func buildEnvironment(env []string) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for _, s := range env {
		// if value is empty, s is like "K=", not "K".
		if !strings.Contains(s, "=") {
			return result, fmt.Errorf("unexpected environment %q", s)
		}
		kv := strings.SplitN(s, "=", 2)
		result[kv[0]] = kv[1]
	}
	return result, nil
}

func GetServicesDeclaredNetworks(serviceConfigs []composetypes.ServiceConfig) map[string]struct{} {
	serviceNetworks := map[string]struct{}{}
	for _, serviceConfig := range serviceConfigs {
		if len(serviceConfig.Networks) == 0 {
			serviceNetworks["default"] = struct{}{}
			continue
		}
		for network := range serviceConfig.Networks {
			serviceNetworks[network] = struct{}{}
		}
	}
	return serviceNetworks
}
