package store

import "eastwh/internal/model"

type UserRepository interface {
	Add(model.User) (model.User, error)
	Login(string, string) (model.User, error)
	Logout(uint) error
	Restore(string) (string, error)
	ChangePassword(uint, string) error
	All() ([]model.User, error)
	Profile(uint) (model.User, error)
	Update(model.User) (model.User, error)
	ByID(uint) (model.User, error)
	ByEmail(string) (model.User, error)
	UpdateToken(uint, string) error
	BlockedUser(uint, bool) error
}
