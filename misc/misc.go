package misc

import (
	"os"

	"github.com/cuigh/auxo/config"
	"github.com/cuigh/auxo/data"
)

const (
	// Version is the version of Swirl
	Version = "0.5.2"
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
	options := config.App().GetSection("swirl")
	DockerHost = loadOption(options, "docker_endpoint", envDockerEndpoint)
	DBType = loadOption(options, "db_type", envDBType)
	DBAddress = loadOption(options, "db_address", envDBAddress)
}

func loadOption(options data.Options, key, env string) (opt string) {
	if options != nil {
		opt = options.String(key)
	}
	if opt == "" {
		opt = os.Getenv(env)
	}
	return
}
