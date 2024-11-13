package store

import "eastwh/internal/model"

type EmployeeTeamRepository interface {
	Add(model.EmployeeTeam) (model.EmployeeTeam, error)
	All() ([]model.EmployeeTeam, error)
	ByID(uint) (model.EmployeeTeam, error)
	ByEmployeeID(uint) ([]model.EmployeeTeam, error)
	ByTeamID(uint) ([]model.EmployeeTeam, error)
	Update(model.EmployeeTeam) (model.EmployeeTeam, error)
	Delete(uint) error
	DeleteEmployeeTeam(uint, uint) error
}
