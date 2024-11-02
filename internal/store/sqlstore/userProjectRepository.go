package sqlstore

import "eastwh/internal/model"

type UserProjectRepository struct {
	store *Store
}

func (r *UserProjectRepository) Add(u model.UserProject) (model.UserProject, error) {
	return u, r.store.db.Create(&u).Error
}

func (r *UserProjectRepository) Update(userproject model.UserProject) (model.UserProject, error) {
	return userproject, r.store.db.Save(&userproject).Error
}

func (r *UserProjectRepository) Delete(id uint) error {
	var up model.UserProject
	result := r.store.db.Table("userprojects").Where("id=,", id)
	err := result.First(&up)
	if err != nil {
		return nil
	}
	return r.store.db.Delete(&up).Error
}

func (r *UserProjectRepository) ByID(ID uint) (up model.UserProject, err error) {
	return up, r.store.db.First(&up, ID).Error
}

func (r *UserProjectRepository) ByUserID(userID uint) (userproject []model.UserProject, err error) {
	return userproject, r.store.db.Where("user_id").Find(&userproject).Error
}

func (r *UserProjectRepository) ByProjectID(projectID uint) (userproject []model.UserProject, err error) {
	return userproject, r.store.db.Where("project_id", projectID).Find(&userproject).Error
}

func (r *UserProjectRepository) All() (userproject []model.UserProject, err error) {
	return userproject, r.store.db.Find(&userproject).Error
}
