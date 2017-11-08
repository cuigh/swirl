package model

import "time"

type Archive struct {
	ID        string    `bson:"_id" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name,omitempty"`
	Content   string    `bson:"content" json:"content,omitempty"`
	CreatedBy string    `bson:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
	UpdatedBy string    `bson:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}

type ArchiveListArgs struct {
	Name      string `bind:"name"`
	PageIndex int    `bind:"page"`
	PageSize  int    `bind:"size"`
}

type Template struct {
	ID        string    `bson:"_id" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name,omitempty"`
	Content   string    `bson:"content" json:"content,omitempty"`
	CreatedBy string    `bson:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
	UpdatedBy string    `bson:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}

type TemplateListArgs struct {
	Name      string `bind:"name"`
	PageIndex int    `bind:"page"`
	PageSize  int    `bind:"size"`
}
