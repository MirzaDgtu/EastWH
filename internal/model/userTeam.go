package model

import "gorm.io/gorm"

type UserTeam struct {
	gorm.Model
	TeamID uint `gorm:"column:team_id" json:"team_id"`
	UserID uint `gorm:"column:user_id" json:"user_id"`
}

func (UserTeam) TableName() string {
	return "user_teams"
}
