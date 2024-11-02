package store

import "eastwh/internal/model"

type ProjectRepository interface {
	Add(model.Project) (model.Project, error)
	ByID(uint) (model.Project, error)
	All() ([]model.Project, error)
	Update(model.Project) (model.Project, error)
	Delete(uint) error
}