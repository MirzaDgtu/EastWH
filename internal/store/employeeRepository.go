package store

import "eastwh/internal/model"

type EmployeeRepository interface {
	Add(model.Employee) (model.Employee, error)
	All() ([]model.Employee, error)
	ByID(uint) (model.Employee, error)
	ByCode(string) (model.Employee, error)
	Update(model.Employee) (model.Employee, error)
	Delete(uint) error
}