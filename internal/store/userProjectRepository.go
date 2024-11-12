package store

import "eastwh/internal/model"

type UserProjectRepository interface {
	Add(model.UserProject) (model.UserProject, error)
	All() ([]model.UserProject, error)
	ByID(uint) (model.UserProject, error)
	ByUserID(uint) ([]model.UserProject, error)
	ByProjectID(uint) ([]model.UserProject, error)
	Update(model.UserProject) (model.UserProject, error)
	Delete(uint) error
	DeleteUserProject(uint, uint) error
}
