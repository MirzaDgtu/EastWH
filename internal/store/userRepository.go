package store

import "eastwh/internal/model"

type UserRepository interface {
	Add(model.User) (model.User, error)
	Login(string, string) (model.User, error)
	Logout(int) error
	Restore(string) (string, error)
	ChangePassword(int, string) error
	All() (users []model.User, error)
	Profile(uint) (model.User, error)
	Update(model.User) (model.User, error)
}
