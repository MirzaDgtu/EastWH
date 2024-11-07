package sqlstore

import "eastwh/internal/model"

type RoleRepository struct {
	store *Store
}

func (r *RoleRepository) Add(u model.Role) (model.Role, error) {
	return u, r.store.db.Create(&u).Error
}

func (r *RoleRepository) All() (roles []model.Role, err error) {
	return roles, r.store.db.Find(&roles).Error
}

func (r *RoleRepository) ByID(id uint) (role model.Role, err error) {
	return role, r.store.db.First(&role, id).Error
}

func (r *RoleRepository) Update(model.Role) (role model.Role, err error) {
	return role, r.store.db.Table("roles").Save(&role).Error
}

func (r *RoleRepository) Delete(id uint) error{
	var role model.Role
	result := r.store.db.Table("roles").Where("id=?", id)
	err := result.First(&role).Error
	if err!= nil{
		return err 
	}
	return r.store.db.Delete(&role).Error
}

