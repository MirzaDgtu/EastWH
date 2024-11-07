package model

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderUid      int     `gorm:"column:order_uid" json:"order_uid"`
	UnicumNum     int     `gorm:"column:unicum_num" json:"unicum_num"`
	OrderDate     string  `gorm:"column:order_date" json:"order_date"`
	OrderSum      float64 `gorm:"column:order_sum" json:"order_sum"`
	Driver        string  `gorm:"column:driver" json:"driver"`
	Agent         string  `gorm:"column:agent" json:"agent"`
	Brieforg      string  `gorm:"column:brieforg" json:"brieforg"`
	ClientId      int     `gorm:"column:client_id" json:"client_id"`
	ClientName    string  `gorm:"column:client_name" json:"client_name"`
	ClientAddress string  `gorm:"column:client_address" json:"client_address"`
	VidDoc        string  `gorm:"column:vid_doc" json:"vid_doc"`
	StartAt       string  `gorm:"column:start_at" json:"start_at"`
	FinishAt      string  `gorm:"column:finish_at" json:"finish_at"`
	Done          bool    `gorm:"column:done" json:"done"`
	Status        int     `gorm:"column:status" json:"status"`
	UserID        uint    `gorm:"column:user_id" json:"user_id"`
	CollectorID   uint    `gorm:"column:collector_id" json:"collector_id"`
}

func (Order) TableName() string {
	return "orders"
}
