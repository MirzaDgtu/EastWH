package sqlstore

import "eastwh/internal/model"

type OrderRepository struct {
	store *Store
}

func (r *OrderRepository) Add(u model.Order) (model.Order, error) {
	err := r.store.db.Create(&u).Error
	return u, err
}

func (r *OrderRepository) Collector(orderid uint, keeper_id uint, collector_id uint) (bool, error) {
	return true, nil
}

func (r *OrderRepository) ByUserID(UserId uint) (order []model.Order, err error) {
	return order, r.store.db.Where("orders.user_id=?", UserId).Find(&order).Error
}

func (r *OrderRepository) ByDateRange(dtStart string, dtFinish string) (orders []model.Order, err error) {
	return orders, r.store.db.Where("orders.order_date BETWEEN ? AND ?", dtStart, dtFinish).Find(&orders).Error
}

func (r *OrderRepository) All() (orders []model.Order, err error) {
	return orders, r.store.db.Find(&orders).Error
}
