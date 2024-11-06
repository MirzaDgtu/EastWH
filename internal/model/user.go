package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string    `gorm:"column:first_name" json:"first_name"`
	Name      string    `gorm:"column:"name" json:"name" validate:"required"`
	LastName  string    `gorm:"column:last_name" json:"last_name" validate:"required"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password,omitempty" validate:"required"`
	LoggedIn  bool      `gorm:"column:loggedin" json:"loggedin"`
	Token     string    `gorm:"column:token" json:"token,omitempty"`
	Restore   bool      `gorm:"column:restore" json:"restore"`
	Blocked   bool      `gorm:"column:blocked" json:"blocked"`
	Phone     string    `gorm:"column:phone" json:"phone"`
	Teams     []Team    `gorm:"many2many:user_teams;" json:"teams"`
	Roles     []Role    `gorm:"many2many:user_roles;" json:"user_roles"`
	Projects  []Project `gorm:"many2many:user_projects;" json:"user_projects"`
}

func (User) TableName() string {
	return "users"
}
