package store

import "eastwh/internal/model"

type RoleRepository interface {
	Add(model.Role) (model.Role, error)
	All() ([]model.Role, error)
	ByID(uint) (model.Role, error)
	Update(model.Role) (model.Role, error)
	Delete(uint) error
}
