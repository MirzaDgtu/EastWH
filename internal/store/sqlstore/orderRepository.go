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
	return orders, r.store.db.Find(&orders).Error
}

func (r *OrderRepository) ByAccessUser(userID uint, startDT, finishDT string) (orders []model.Order, err error) {
	return orders, r.store.db.Raw(`
								SELECT o.*
								FROM user_projects usp
									LEFT JOIN projects p on usp.project_id = p.id
									LEFT JOIN users u on u.id = usp.user_id
									LEFT JOIN orders o on o.vid_doc	= p.vid_doc
								WHERE o.deleted_at IS NULL
									  AND usp.deleted_at IS NULL
									  AND p.deleted_at IS NULL
									  AND	usp.user_id = ?
									  AND o.folio_date between ? and ?
									  AND o.Done = 0								
									  AND o.Check = 0
`, userID, startDT, finishDT).Scan(&orders).Error
}

func (r *OrderRepository) AssemblyOrder(startDT, finishDT string) (assemblyOrders []model.AssemblyOrder, err error) {
	/*	return assemblyOrders, r.store.db.Raw(`
												SELECT *
												FROM eastwh.assembly_orders_vw ao
												where ao.folio_date between ? and ?
		`, startDT, finishDT).Scan(&assemblyOrders).Error
	*/
	return assemblyOrders, r.store.db.Raw("CALL `eastwh`.`get_assembly_orders`(?, ?)", startDT, finishDT).Scan(&assemblyOrders).Error
}

func (r *OrderRepository) SetCheck(orderuid uint, user_id uint, check bool) error {
	err := r.store.db.Model(&model.Order{}).Where("order_uid=?", orderuid).Updates(map[string]interface{}{
		"user_id": user_id,
		"check":   check,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) CheckedList(startDT, finishDT string, checkStatus bool) (orders []model.Order, err error) {
	return orders, r.store.db.Raw(`SELECT * 
											FROM eastwh.orders o
											where o.folio_date between ? and ?
											AND IFNULL(o.check, 0) = ?
	`, startDT, finishDT, checkStatus).Scan(&orders).Error
}
