package model

import "gorm.io/gorm"

type EmployeeTeam struct {
	gorm.Model
	TeamID     uint `gorm:"primaryKey;autoIncrement:false" json:"team_id"`
	EmployeeID uint `gorm:"primaryKey;autoIncrement:false" json:"employee_id"`
}

func (EmployeeTeam) TableName() string {
	return "employee_teams"
}
