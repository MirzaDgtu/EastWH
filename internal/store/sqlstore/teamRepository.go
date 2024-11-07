package sqlstore

import "eastwh/internal/model"

type TeamRepository struct {
	store *Store
}

func (r *TeamRepository) Add(u model.Team) (model.Team, error) {
	return u, r.store.db.Create(&u).Error
}

func (r *TeamRepository) ByID(id uint) (team model.Team, err error) {
	return team, r.store.db.Where("id=?", id).First(&team).Error
}

func (r *TeamRepository) All() (teams []model.Team, err error) {
	return teams, r.store.db.Model(&model.Team{}).Find(&teams).Error
}

func (r *TeamRepository) Update(u model.Team) (model.Team, error) {
	return u, r.store.db.Model(&u).Update("name", u.Name).Error
}

func (r *TeamRepository) Delete(id uint) error {
	var team model.Team
	result := r.store.db.Table("teams").Where("id=?", id)
	err := result.First(&team).Error
	if err != nil {
		return err
	}
	return r.store.db.Delete(&team).Error
}
