package misc

import (
	"time"

	"github.com/cuigh/auxo/config"
)

const (
	keyDockerEndpoint = "swirl.docker_endpoint"
	keyDBType         = "swirl.db_type"
	keyDBAddress      = "swirl.db_address"
	keyAuthTimeout    = "swirl.auth_timeout"
	envDockerEndpoint = "DOCKER_ENDPOINT"
	envDBType         = "DB_TYPE"
	envDBAddress      = "DB_ADDRESS"
	envAuthTimeout    = "AUTH_TIMEOUT"
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

// Options holds custom options of swirl.
var Options = &struct {
	DockerEndpoint string
	DBType         string
	DBAddress      string
	AuthTimeout    time.Duration
}{
	DBType:      "mongo",
	DBAddress:   "localhost:27017/swirl",
	AuthTimeout: 30 * time.Minute,
}

// BindOptions binds options to environment variables.
func BindOptions() {
	config.BindEnv(keyDockerEndpoint, envDockerEndpoint)
	config.BindEnv(keyDBType, envDBType)
	config.BindEnv(keyDBAddress, envDBAddress)
	config.BindEnv(keyAuthTimeout, envAuthTimeout)
}
