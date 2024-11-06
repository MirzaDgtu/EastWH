package sqlstore

import "eastwh/internal/model"

type EmployeeRepository struct {
	store *Store
}

func (r *EmployeeRepository) Add(u model.Employee) (model.Employee, error) {
	return u, r.store.db.Create(&u).Error
}

func (r *EmployeeRepository) All() (employee []model.Employee, err error) {
	return employee, r.store.db.Table("employees").Select("*").Scan(&employee).Error
}

func (r *EmployeeRepository) ByID(id uint) (employee model.Employee, err error) {
	employee.ID = id
	return employee, r.store.db.First(&employee).Error
}

func (r *EmployeeRepository) ByCode(code string) (employee model.Employee, err error) {
	employee.Code = code
	return employee, r.store.db.First(&employee).Error
}

func (r *EmployeeRepository) Update(u model.Employee) (model.Employee, error) {
	return u, r.store.db.Model(&model.Employee{}).Where("id=?", u.ID).Updates(map[string]interface{}{
		"code":       u.Code,
		"first_name": u.FirstName,
		"name":       u.Name,
		"last_name":  u.LastName,
		"inn":        u.INN,
		"phone":      u.Phone,
	}).Error
}

func (r *EmployeeRepository) Delete(id uint) error {
	var employee model.Employee
	result := r.store.db.Table("employees").Where("id=?", id)
	err := result.First(&employee).Error
	if err != nil {
		return err
	}
	return r.store.db.Delete(&employee).Error
}
