package misc

import (
	"os"

	"github.com/cuigh/auxo/config"
)

const (
	// Version is the version of Swirl
	Version = "0.6"
)

const (
	envDockerEndpoint = "DOCKER_ENDPOINT"
	envDBType         = "DB_TYPE"
	envDBAddress      = "DB_ADDRESS"
)

var (
	DockerHost string
	DBType     string
	DBAddress  string
)

func init() {
	DockerHost = loadOption("swirl.docker_endpoint", envDockerEndpoint)
	DBType = loadOption("swirl.db_type", envDBType)
	DBAddress = loadOption("swirl.db_address", envDBAddress)
}

func loadOption(key, env string) (opt string) {
	opt = config.GetString(key)
	if opt == "" {
		opt = os.Getenv(env)
	}
	return
}
