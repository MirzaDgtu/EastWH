package sqlstore

import (
	"crypto/rand"
	"eastwh/internal/model"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Add(u model.User) (model.User, error) {
	hashPassword(&u.Password)
	err := r.store.db.Create(&u).Error
	u.Password = ""
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

	user.Password = ""
	return user, nil
}

func checkPassword(existingHash, incomingPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}

func (r *UserRepository) Logout(id uint) error {
	user := model.User{
		Model: gorm.Model{
			ID: id,
		},
	}

	return r.store.db.Model(&user).Where("id = ?", id).Updates(map[string]interface{}{"loggedin": 0,
		"token": ""}).Error
}

func (r *UserRepository) UpdateToken(id uint, token string) error {
	return r.store.db.Model(&model.User{}).Where("id=?", id).Update("token", token).Error
}

func (r *UserRepository) Restore(email string) (password string, err error) {
	pass, err := generateTemporaryPassword()
	if err != nil {
		return "", err
	}

	var user model.User
	result := r.store.db.Table("users").Where("email=?", email)
	err = result.First(&user).Error
	if err != nil {
		return "", err
	}

	hPass := pass
	hashPassword(&hPass)
	return pass, r.store.db.Model(&model.User{}).Where("id=?", user.ID).Updates(map[string]interface{}{"password": hPass,
		"restore": true}).Error
}

func (r *UserRepository) ChangePassword(id uint, password string) error {
	err := hashPassword(&password)
	if err != nil {
		return err
	}

	fmt.Println("id - ", id)
	return r.store.db.Model(&model.User{}).Where("id=?", id).Updates(map[string]interface{}{"password": password,
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

func (r *UserRepository) Profile(id uint) (u model.User, err error) {
	err = r.store.db.
		//	Preload("TeamUsers.Team").
		//Preload("TeamUsers.Employee").
		Preload("Teams"). // если нужны и сами команды
		//Preload("Teams.Employees"). // если нужны сотрудники команд
		First(&u, id).Error
	if err != nil {
		return model.User{}, err
	}
	u.Password = ""
	return u, nil
}

func (r *UserRepository) Update(u model.User) (model.User, error) {
	err := r.store.db.Model(&u).Updates(map[string]interface{}{"first_name": u.FirstName,
		"last_name": u.LastName,
		"name":      u.Name,
		"phone":     u.Phone}).Error
	if err != nil {
		return model.User{}, err
	}

	u.Password = ""
	return u, nil
}

func (r *UserRepository) ByID(id uint) (u model.User, err error) {
	return u, r.store.db.First(&u, id).Error
}

func (r *UserRepository) ByEmail(email string) (u model.User, err error) {
	return u, r.store.db.Where("email=?", email).First(&u).Error
}

func (r *UserRepository) BlockedUser(id uint, blocked bool) error {
	return r.store.db.Model(&model.User{}).Where("id=?", id).Update("blocked", blocked).Error
}

// Функция для генерации временного пароля
func generateTemporaryPassword() (string, error) {
	bytes := make([]byte, 6) // Длина пароля
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
