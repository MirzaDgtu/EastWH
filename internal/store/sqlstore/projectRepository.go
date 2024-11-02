package sqlstore
<<<<<<< HEAD
=======

import "eastwh/internal/model"

type ProjectRepository struct {
	store *Store
}

func (r *ProjectRepository) Add(u model.Project) (model.Project, error) {
	return u, r.store.db.Create(&u).Error
}

func (r *ProjectRepository) All() (project []model.Project, err error) {
	return project, r.store.db.Find(&project).Error
}

func (r *ProjectRepository) ByID(id uint) (project model.Project, err error) {
	return project, r.store.db.First(&project, id).Error
}

func (r *ProjectRepository) Update(model.Project) (project model.Project, err error) {
	return project, r.store.db.Table("projects").Save(&project).Error
}

func (r *ProjectRepository) Delete(id uint) error {
	var project model.Project
	result := r.store.db.Table("projects").Where("id=?", id)
	err := result.First(&project).Error
	if err != nil {
		return err
	}
	return r.store.db.Delete(&project).Error
}
>>>>>>> 7806773d89dd16394f5140f59697b6e30d089ac6
