package store

import "eastwh/internal/model"

type OrderRepository interface {
	Add(model.Order) (model.Order, error)
	SetCollector(orderuid uint, user_id uint, collector_id uint) error
	ByUserID(uint) ([]model.Order, error)
	ByID(uint) ([]model.Order, error)
	ByOrderUID(uint) ([]model.Order, error)
	ByDateRange(string, string) ([]model.Order, error)
	All() ([]model.Order, error)
}
