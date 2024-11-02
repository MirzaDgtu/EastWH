package store

import "eastwh/internal/model"

type OrderRepository interface {
	Add([]model.Order) (model.Order)
	Collector(orderid uint,keeper_id uint,collector_id uint) (bool, error) 
	ByUserID(uint) ([]model.Order, error) 
	ByDateRange(string, string) ([]model.Order, error)
	All() ([]model.Order)
}