package model

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name   string `gorm:"column:name;not null;unique" json:"name" validate:"required"`
	VidDoc string `json:"vid_doc"`
	Users  []User `gorm:"many2many:user_projects;" json:"users"`
}

func (Project) TableName() string {
	return "projects"
}
