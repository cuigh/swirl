package model

import "time"

type Setting struct {
	LDAP struct {
		Enabled   bool   `bson:"enabled" json:"enabled,omitempty"`
		Address   string `bson:"address" json:"address,omitempty"`
		BaseDN    string `bson:"base_dn" json:"base_dn,omitempty"`
		NameAttr  string `bson:"name_attr" json:"name_attr,omitempty"`
		EmailAttr string `bson:"email_attr" json:"email_attr,omitempty"`
	} `bson:"ldap" json:"ldap,omitempty"`
	TimeZone struct {
		Name   string `bson:"name" json:"name,omitempty"`     // Asia/Shanghai
		Offset int32  `bson:"offset" json:"offset,omitempty"` // seconds east of UTC
	} `bson:"tz" json:"tz,omitempty"`
	Language  string    `bson:"lang" json:"lang,omitempty"`
	UpdatedBy string    `bson:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}
