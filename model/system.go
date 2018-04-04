package model

import "time"

// Perm control scope
const (
	PermNone      = 0
	PermWrite     = 1
	PermReadWrite = 2
)

// Setting represents the options of swirl.
type Setting struct {
	Version string
	LDAP    struct {
		Enabled  bool   `bson:"enabled" json:"enabled,omitempty"`
		Address  string `bson:"address" json:"address,omitempty"`
		Security int32  `bson:"security" json:"security,omitempty"` // 0-None/1-TLS/2-StartTLS
		//TLSCert        string `bson:"tls_cert" json:"tls_cert,omitempty"`       // TLS cert
		//TLSVerify      bool   `bson:"tls_verify" json:"tls_verify,omitempty"`   // Verify cert
		Authentication int32  `bson:"auth" json:"auth,omitempty"`               // 0-Simple/1-Bind
		BindDN         string `bson:"bind_dn" json:"bind_dn,omitempty"`         // DN to bind with
		BindPassword   string `bson:"bind_pwd" json:"bind_pwd,omitempty"`       // Bind DN password
		BaseDN         string `bson:"base_dn" json:"base_dn,omitempty"`         // Base search path for users
		UserDN         string `bson:"user_dn" json:"user_dn,omitempty"`         // Template for the DN of the user for simple auth
		UserFilter     string `bson:"user_filter" json:"user_filter,omitempty"` // Search filter for user
		NameAttr       string `bson:"name_attr" json:"name_attr,omitempty"`
		EmailAttr      string `bson:"email_attr" json:"email_attr,omitempty"`
	} `bson:"ldap" json:"ldap,omitempty"`
	TimeZone struct {
		Name   string `bson:"name" json:"name,omitempty"`     // Asia/Shanghai
		Offset int32  `bson:"offset" json:"offset,omitempty"` // seconds east of UTC
	} `bson:"tz" json:"tz,omitempty"`
	Language string `bson:"lang" json:"lang,omitempty"`
	Metrics  struct {
		Prometheus string `bson:"prometheus" json:"prometheus"`
	} `bson:"metrics" json:"metrics"`
	UpdatedBy string    `bson:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}

// Perm holds permissions of Docker resource.
type Perm struct {
	ResType string   `json:"res_type"`
	ResID   string   `json:"res_id"`
	Scope   int32    `json:"scope"`
	Roles   []string `json:"roles"`
	Users   []string `json:"users"`
}
