package model

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Code      string `json:"code"`
	FirstName string `json:"first_name" validate:"required"`
	Name      string `json:"name" validate:"required"`
	LastName  string `json:"last_name"`
	INN       string `gorm:"column:inn" json:"inn"`
	Phone     string `json:"phone"`
	Owner string `gorm:"column:owner" json:"owner"`
}

func (Employee) TableName() string {
	return "employees"
}
