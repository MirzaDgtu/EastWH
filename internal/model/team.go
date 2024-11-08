package model

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name      string     `gorm:"column:name;not null;unique" json:"name" validate:"required"`
	Users     []User     `gorm:"many2many:user_teams;" json:"users"`
	Employees []Employee `gorm:"many2many:user_teams" json:"employees"`
	TeamUsers []UserTeam `gorm:"foreignKey:TeamID" json:"team_users,omitempty"`
}

func (Team) TableName() string {
	return "teams"
}
