package store

import "eastwh/internal/model"

type UserRoleRepository interface {
	Add(model.UserRole) (model.UserRole, error)
	All() ([]model.UserRole, error)
	ByID(uint) (model.UserRole, error)
	ByUserID(uint) ([]model.UserRole, error)
	ByRoleID(uint) ([]model.UserRole, error)
	Update(model.UserRole) (model.UserRole, error)
	Delete(uint) error
}
