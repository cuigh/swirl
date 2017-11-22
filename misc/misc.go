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

// TimeZones holds some commonly used time-zones.
var TimeZones = []struct {
	Name   string
	Offset int32 // seconds east of UTC
}{
	{"GMT", 0},
	{"GMT+12", 12 * 60 * 60},
	{"GMT+11", 11 * 60 * 60},
	{"GMT+10", 10 * 60 * 60},
	{"GMT+9", 9 * 60 * 60},
	{"GMT+8(Asia/Shanghai)", 8 * 60 * 60},
	{"GMT+7", 7 * 60 * 60},
	{"GMT+6", 6 * 60 * 60},
	{"GMT+5", 5 * 60 * 60},
	{"GMT+4", 4 * 60 * 60},
	{"GMT+3", 3 * 60 * 60},
	{"GMT+2", 2 * 60 * 60},
	{"GMT+1", 1 * 60 * 60},
	{"GMT-1", -1 * 60 * 60},
	{"GMT-2", -2 * 60 * 60},
	{"GMT-3", -3 * 60 * 60},
	{"GMT-4", -4 * 60 * 60},
	{"GMT-5", -5 * 60 * 60},
	{"GMT-6", -6 * 60 * 60},
	{"GMT-7", -7 * 60 * 60},
	{"GMT-8", -8 * 60 * 60},
	{"GMT-9", -9 * 60 * 60},
	{"GMT-10", -10 * 60 * 60},
	{"GMT-11", -11 * 60 * 60},
	{"GMT-12", -12 * 60 * 60},
}

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
