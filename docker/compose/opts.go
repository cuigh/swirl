package compose

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
)

// ConvertKVStringsToMapWithNil converts ["key=value"] to {"key":"value"}
// but set unset keys to nil - meaning the ones with no "=" in them.
// We use this in cases where we need to distinguish between
//   FOO=  and FOO
// where the latter case just means FOO was mentioned but not given a value
func ConvertKVStringsToMapWithNil(values []string) map[string]*string {
	result := make(map[string]*string, len(values))
	for _, value := range values {
		kv := strings.SplitN(value, "=", 2)
		if len(kv) == 1 {
			result[kv[0]] = nil
		} else {
			result[kv[0]] = &kv[1]
		}
	}

	return result
}

// ParseRestartPolicy returns the parsed policy or an error indicating what is incorrect
func ParseRestartPolicy(policy string) (container.RestartPolicy, error) {
	p := container.RestartPolicy{}

	if policy == "" {
		return p, nil
	}

	parts := strings.Split(policy, ":")

	if len(parts) > 2 {
		return p, fmt.Errorf("invalid restart policy format")
	}
	if len(parts) == 2 {
		count, err := strconv.Atoi(parts[1])
		if err != nil {
			return p, fmt.Errorf("maximum retry count must be an integer")
		}

		p.MaximumRetryCount = count
	}

	p.Name = parts[0]

	return p, nil
}

// ParseCPUs takes a string ratio and returns an integer value of nano cpus
func ParseCPUs(value string) (int64, error) {
	cpu, ok := new(big.Rat).SetString(value)
	if !ok {
		return 0, fmt.Errorf("failed to parse %v as a rational number", value)
	}
	nano := cpu.Mul(cpu, big.NewRat(1e9, 1))
	if !nano.IsInt() {
		return 0, fmt.Errorf("value is too precise")
	}
	return nano.Num().Int64(), nil
}
