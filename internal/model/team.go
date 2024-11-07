package model

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name  string `gorm:"column:name;not null;unique" json:"name" validate:"required"`
	Users []User `gorm:"many2many:user_teams;" json:"users"`
}

func (Team) TableName() string {
	return "teams"
}
