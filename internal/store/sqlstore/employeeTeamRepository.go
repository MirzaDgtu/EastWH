package sqlstore

import "eastwh/internal/model"

type EmployeeTeamRepository struct {
	store *Store
}

func (r *EmployeeTeamRepository) Add(u model.EmployeeTeam) (model.EmployeeTeam, error) {
	return u, r.store.db.Create(&u).Error
}

func (r *EmployeeTeamRepository) All() (et []model.EmployeeTeam, err error) {
	return et, r.store.db.Find(&et).Error
}

func (r *EmployeeTeamRepository) Update(et model.EmployeeTeam) (model.EmployeeTeam, error) {
	return et, r.store.db.Save(et).Error
}

func (r *EmployeeTeamRepository) Delete(id uint) error {
	return r.store.db.Exec(`DELETE FROM employee_teams WHERE id=?`, id).Error
}

func (r *EmployeeTeamRepository) DeleteEmployeeTeam(employee_id, team_id uint) error {
	return r.store.db.
		Exec("DELETE FROM employee_teams WHERE team_id = ? and employee_id = ?",
			team_id, employee_id).Error
}

func (r *EmployeeTeamRepository) ByID(id uint) (et model.EmployeeTeam, err error) {
	return et, r.store.db.First(&et, id).Error
}

func (r *EmployeeTeamRepository) ByEmployeeID(employeeID uint) (et []model.EmployeeTeam, err error) {
	return et, r.store.db.Where("employee_id = ?", employeeID).Find(&et).Error
}

func (r *EmployeeTeamRepository) ByTeamID(teamID uint) (et []model.EmployeeTeam, err error) {
	return et, r.store.db.Where("team_id = ?", teamID).Find(&et).Error
}
