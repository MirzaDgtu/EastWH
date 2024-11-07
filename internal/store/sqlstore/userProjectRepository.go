package sqlstore

import "eastwh/internal/model"

type UserProjectRepository struct {
	store *Store
}

func (r *UserProjectRepository) Add(u model.UserProject) (model.UserProject, error) {
	return u, r.store.db.Create(&u).Error
}

func (r *UserProjectRepository) Update(userproject model.UserProject) (model.UserProject, error) {
	return userproject, r.store.db.Model(&userproject).Update("user_id", userproject.UserID).Update("project_id", userproject.ProjectID).Error
}

func (r *UserProjectRepository) Delete(id uint) error {
	var up model.UserProject
	result := r.store.db.Table("user_projects").Where("id=?", id)
	err := result.First(&up)
	if err != nil {
		return nil
	}
	return r.store.db.Delete(&up).Error
}

func (r *UserProjectRepository) ByID(Id uint) (up model.UserProject, err error) {
	up.ID = Id
	return up, r.store.db.First(&up, Id).Error
}

func (r *UserProjectRepository) ByUserID(userID uint) (userproject []model.UserProject, err error) {
	return userproject, r.store.db.Where("user_id=?", userID).Find(&userproject).Error
}

func (r *UserProjectRepository) ByProjectID(projectID uint) (userproject []model.UserProject, err error) {
	return userproject, r.store.db.Where("project_id=?", projectID).Find(&userproject).Error
}

func (r *UserProjectRepository) All() (userproject []model.UserProject, err error) {
	return userproject, r.store.db.Find(&userproject).Error
}
