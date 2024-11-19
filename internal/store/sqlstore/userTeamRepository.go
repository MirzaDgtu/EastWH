package sqlstore

import "eastwh/internal/model"

type UserTeamRepository struct {
	store *Store
}

func (r *UserTeamRepository) Add(u model.UserTeam) (model.UserTeam, error) {
	return u, r.store.db.Create(&u).Error
}

func (r *UserTeamRepository) Update(userteam model.UserTeam) (model.UserTeam, error) {
	return userteam, r.store.db.Save(&userteam).Error
}

func (r *UserTeamRepository) Delete(id uint) error {
	var ut model.UserTeam
	result := r.store.db.Table("userteams").Where("id=,", id)
	err := result.First(&ut)
	if err != nil {
		return nil
	}
	return r.store.db.Delete(&ut).Error
}

func (r *UserTeamRepository) DeleteUserTeam(team_id, user_id uint) error {
	return r.store.db.Exec("DELETE FROM user_teams WHERE team_id = ? AND user_id = ?", team_id, user_id).Error
}

func (r *UserTeamRepository) ByID(ID uint) (ur model.UserTeam, err error) {
	return ur, r.store.db.First(&ur, ID).Error
}

func (r *UserTeamRepository) ByUserID(userID uint) (userteam []model.UserTeam, err error) {
	return userteam, r.store.db.Where("user_id").Find(&userteam).Error
}

func (r *UserTeamRepository) ByTeamID(teamID uint) (userteam []model.UserTeam, err error) {
	return userteam, r.store.db.Where("team_id", teamID).Find(&userteam).Error
}

func (r *UserTeamRepository) All() (userteam []model.UserTeam, err error) {
	return userteam, r.store.db.Find(&userteam).Error
}
