package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"unique" json:"name"`
	Description string `json:"description" validation:"required"`
	Priority    int    `json:"priority"`
	Users       []User `gorm:"many2many:user_roles;" json:"users"`
}

func (Role) TableName() string {
	return "roles"
}
