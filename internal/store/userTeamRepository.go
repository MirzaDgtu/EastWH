package store

import "eastwh/internal/model"

type UserTeamRepository interface {
	Add(model.UserTeam) (model.UserTeam, error)
	All() ([]model.UserTeam, error)
	ByID(uint) (model.UserTeam, error)
	ByUserID(uint) ([]model.UserTeam, error)
	ByTeamID(uint) ([]model.UserTeam, error)
	Update(model.UserTeam) (model.UserTeam, error)
	Delete(uint) error
}
