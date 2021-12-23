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
	Agents           []string
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
		"agents",
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
		Version string `json:"version"`
	} `json:"system"`
	LDAP struct {
		Enabled        bool   `json:"enabled"`
		Address        string `json:"address"`
		Security       int32  `json:"security"` // 0, 1, 2
		Authentication string `json:"auth"`     // simple, bind
		BindDN         string `json:"bind_dn"`
		BindPassword   string `json:"bind_pwd"`    // Bind DN password
		BaseDN         string `json:"base_dn"`     // Base search path for users
		UserDN         string `json:"user_dn"`     // Template for the DN of the user for simple auth
		UserFilter     string `json:"user_filter"` // Search filter for user
		NameAttr       string `json:"name_attr"`
		EmailAttr      string `json:"email_attr"`
	} `json:"ldap"`
	Metric struct {
		Prometheus string `json:"prometheus"`
	} `json:"metric"`
}

func init() {
	bindOptions()
}
