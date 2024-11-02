package sqlstore

import "eastwh/internal/model"

type TeamRepository struct {
	store *Store
}

func (r *TeamRepository) Add(u model.Team) (model.Team, error) {
	err := r.store.db.Create(&u).Error
	return u, err
}

func (r *TeamRepository) ByID(id uint) (team model.Team, err error) {
	team.ID = uint(id)
	return team, r.store.db.First(&team).Error
}

func (r *TeamRepository) All() (teams []model.Team, err error) {
	return teams, r.store.db.Table("teams").Select("*").Error
}

func (r *TeamRepository) Update(u model.Team) (model.Team, error) {
	return u, r.store.db.Table("teams").Save(&u).Error
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
