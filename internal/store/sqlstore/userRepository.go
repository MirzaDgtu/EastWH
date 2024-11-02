package sqlstore

import "eastwh/internal/model"

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Add(u model.User) (model.User, error) {
	return u, nil
}

func (r *UserRepository) Login(email, password string) (u model.User, err error) {
	return u, nil
}

func (r *UserRepository) Logout(id int) error {
	return nil
}

func (r *UserRepository) Restore(email string) (Password string, err error) {
	return Password, nil
}

func (r *UserRepository) ChangePassword(int, string) error {
	return nil
}

func (r *UserRepository) All() (users []model.User, err error) {
	return users, r.store.db.Preload("").Preload("").Preload("").Find(&users).Error
}

func (r *UserRepository) Profile(uint) (u model.User, err error) {
	return u, nil
}

func (r *UserRepository) Update(u model.User) (model.User, error) {
	err := r.store.db.Model(&u).Updates(map[string]interface{}{"firstname": u.FirstName,
		"lastName": u.LastName,
		"Name":     u.Name,
		"Phone":    u.Phone}).Error
	if err != nil {
		return u, err
	}

	u.Password = ""
	return u, nil
}
