package sqlstore

import (
	"eastwh/internal/model"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Add(u model.User) (model.User, error) {
	err := r.store.db.Create(&u).Error
	return u, err
}

func (r *UserRepository) Login(email, password string) (user model.User, err error) {
	result := r.store.db.Table("users").Where(&model.User{Email: email})
	err = result.First(&user).Error
	if err != nil {
		return user, err
	}

	if !checkPassword(user.Password, password) {
		return user, errors.New("Invalid password")
	}

	err = result.Update("loggedin", 1).Error
	if err != nil {
		return user, err
	}

	err = result.Find(&user).Error
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func checkPassword(existingHash, incomingPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}

func (r *UserRepository) Logout(id int) error {
	user := model.User{
		Model: gorm.Model{
			ID: uint(id),
		},
	}

	return r.store.db.Model(&user).Where("id = ?", id).Updates(map[string]interface{}{"loggedin": 0,
		"token":         "",
		"refresh_token": ""}).Error
}

func (r *UserRepository) Restore(email string) (password string, err error) {
	return password, nil
}

func (r *UserRepository) ChangePassword(id int, password string) error {
	fmt.Println(password)
	err := hashPassword(&password)
	if err != nil {
		return err
	}

	var user model.User
	return r.store.db.Model(&user).Where("id=?", id).Updates(map[string]interface{}{"pass": password,
		"restore": false}).Error
}

func hashPassword(s *string) error {
	if s == nil {
		return errors.New("Reference provided for hashing password is nil")
	}
	//converd password string to byte slice
	sBytes := []byte(*s)
	//Obtain hashed password
	hashedBytes, err := bcrypt.GenerateFromPassword(sBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//update password string with the hashed version
	*s = string(hashedBytes[:])
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
