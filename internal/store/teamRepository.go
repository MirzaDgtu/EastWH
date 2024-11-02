package store

import "eastwh/internal/model"

type TeamRepository interface {
	Add(model.Team) (model.Team, error)
	ByID(uint) (model.Team, error)
	All() ([]model.Team, error)
	Update(model.Team) (model.Team, error)
	Delete(uint) error
}