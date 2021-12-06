package misc

import (
	"strings"
	"time"

	"github.com/cuigh/auxo/config"
	"github.com/cuigh/auxo/errors"
)

// Options holds custom options of Swirl.
var Options = &struct {
	DockerEndpoint   string
	DockerAPIVersion string
	DBType           string
	DBAddress        string
	TokenKey         string
	TokenExpiry      time.Duration
}{
	DBType:      "mongo",
	DBAddress:   "mongodb://localhost:27017/swirl",
	TokenExpiry: 30 * time.Minute,
}

func bindOptions() {
	var keys = []string{
		"docker_endpoint",
		"docker_api_version",
		"db_type",
		"db_address",
		"token_key",
		"token_expiry",
	}
	for _, key := range keys {
		config.BindEnv("swirl."+key, strings.ToUpper(key))
	}
}

func LoadOptions() (err error) {
	err = config.UnmarshalOption("swirl", &Options)
	if err != nil {
		err = errors.Wrap(err, "failed to load options")
	}
	return
}

// Setting represents the settings of Swirl.
type Setting struct {
	System struct {
		Version string
	}
	Region struct {
		Language string `option:"lang"`
		Timezone int32
	}
	LDAP struct {
		Enabled        bool
		Address        string
		Security       int32  // 0, 1, 2
		Authentication string `option:"auth"` // simple, bind
		BindDN         string
		BindPassword   string `option:"bind_pwd"` // Bind DN password
		BaseDN         string // Base search path for users
		UserDN         string // Template for the DN of the user for simple auth
		UserFilter     string // Search filter for user
		NameAttr       string
		EmailAttr      string
	}
	Metric struct {
		Prometheus string
	}
}

func init() {
	bindOptions()
}
