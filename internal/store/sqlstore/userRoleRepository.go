package sqlstore

import "eastwh/internal/model"

type UserRoleRepository struct {
	store *Store
}

func (r *UserRoleRepository) Add(u model.UserRole) (model.UserRole, error) {
	return u, r.store.db.Create(&u).Error
}

func (r *UserRoleRepository) Update(userrole model.UserRole) (model.UserRole, error) {
	return userrole, r.store.db.Save(&userrole).Error
}

func (r *UserRoleRepository) Delete(id uint) error {
	var ur model.UserRole
	result := r.store.db.Table("userroles").Where("id=,", id)
	err := result.First(&ur)
	if err != nil {
		return nil
	}
	return r.store.db.Delete(&ur).Error
}

func (r *UserRoleRepository) ByID(ID uint) (ur model.UserRole, err error) {
	return ur, r.store.db.First(&ur, ID).Error
}

func (r *UserRoleRepository) ByUserID(userID uint) (userroles []model.UserRole, err error) {
	return userroles, r.store.db.Where("user_id").Find(&userroles).Error
}

func (r *UserRoleRepository) ByRoleID(roleID uint) (userroles []model.UserRole, err error) {
	return userroles, r.store.db.Where("role_id", roleID).Find(&userroles).Error
}

func (r *UserRoleRepository) All() (userroles []model.UserRole, err error) {
	return userroles, r.store.db.Find(&userroles).Error
}
