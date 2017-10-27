package misc

import (
	"github.com/cuigh/auxo/config"
)

const (
	keyDockerEndpoint = "swirl.docker_endpoint"
	keyDBType         = "swirl.db_type"
	keyDBAddress      = "swirl.db_address"
	envDockerEndpoint = "DOCKER_ENDPOINT"
	envDBType         = "DB_TYPE"
	envDBAddress      = "DB_ADDRESS"
)

var Options = &struct {
	DockerEndpoint string
	DBType         string
	DBAddress      string
}{}

func BindOptions() {
	config.BindEnv(keyDockerEndpoint, envDockerEndpoint)
	config.BindEnv(keyDBType, envDBType)
	config.BindEnv(keyDBAddress, envDBAddress)
}

func LoadOptions() {
	Options.DockerEndpoint = config.GetString(keyDockerEndpoint)
	Options.DBType = config.GetString(keyDBType)
	Options.DBAddress = config.GetString(keyDBAddress)
}
