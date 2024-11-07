package sqlstore

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
	project.ID = id
	return project, r.store.db.First(&project, id).Error
}

func (r *ProjectRepository) Update(u model.Project) (project model.Project, err error) {
	return u, r.store.db.Model(&u).Update("name", u.Name).Error
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
