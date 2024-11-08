package model

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Code      string `gorm:"not null;unique" json:"code"`
	FirstName string `json:"first_name" validate:"required"`
	Name      string `json:"name" validate:"required"`
	LastName  string `json:"last_name"`
	INN       string `gorm:"column:inn" json:"inn"`
	Phone     string `json:"phone"`
	Teams     []Team `gorm:"many2many:user_teams" json:"teams,omitempty"`
	Users     []User `gorm:"many2many:user_teams" json:"users,omitempty"`
}

func (Employee) TableName() string {
	return "employees"
}
