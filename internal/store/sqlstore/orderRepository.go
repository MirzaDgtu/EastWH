package sqlstore

import "eastwh/internal/model"

type OrderRepository struct {
	store *Store
}

func (r *OrderRepository) Add(u model.Order) (model.Order, error) {
	return u, r.store.db.Create(&u).Error
}

func (r *OrderRepository) SetCollector(orderuid uint, user_id uint, employee_id uint) error {
	err := r.store.db.Model(&model.Order{}).Where("order_uid=?", orderuid).Updates(map[string]interface{}{
		"user_id":     user_id,
		"employee_id": employee_id,
		"done":        1,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) ByUserID(userID uint) (order []model.Order, err error) {
	return order, r.store.db.Where("user_id=?", userID).Find(&order).Error
}

func (r *OrderRepository) ByID(ID uint) (order []model.Order, err error) {
	return order, r.store.db.Where("id=?", ID).Find(&order).Error
}

func (r *OrderRepository) ByOrderUID(orderUID uint) (order []model.Order, err error) {
	return order, r.store.db.Where("order_uid=?", orderUID).Find(&order).Error
}

func (r *OrderRepository) ByDateRange(dtStart string, dtFinish string) (orders []model.Order, err error) {
	return orders, r.store.db.Where("orders.order_date BETWEEN ? AND ?", dtStart, dtFinish).Find(&orders).Error
}

func (r *OrderRepository) All() (orders []model.Order, err error) {
	return orders, r.store.db.Where("done=?", 0).Find(&orders).Error
}
