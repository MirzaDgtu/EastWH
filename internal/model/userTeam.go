package model

import "gorm.io/gorm"

type UserTeam struct {
	gorm.Model
	TeamID uint `gorm:"column:team_id" json:"team_id"`
	UserID uint `gorm:"column:user_id" json:"user_id"`
	//EmployeeID uint `gorm:"column:employee_id" json:"employee_id"`
	Team Team `gorm:"foreignKey:TeamID" json:"team"`
	User User `gorm:"foreignKey:UserID" json:"user"`
	//	Employee   Employee `gorm:"foreignKey:EmployeeID" json:"employee"`
}

func (UserTeam) TableName() string {
	return "user_teams"
}
