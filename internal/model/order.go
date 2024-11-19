package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderUid      int     `gorm:"column:order_uid;not null;unique" json:"order_uid"`
	UnicumNum     int     `gorm:"column:unicum_num" json:"unicum_num"`
	FolioNum      int     `gorm:"column:folio_num" json:"folio_num"`
	FolioDate     string  `gorm:"column:folio_date" json:"folio_date"`
	OrderDate     string  `gorm:"column:order_date;size:19" json:"order_date"`
	OrderSum      float64 `gorm:"column:order_sum" json:"order_sum"`
	FolioSum      float64 `gorm:"column:folio_sum" json:"folio_sum"`
	Driver        string  `gorm:"column:driver;size:100" json:"driver"`
	Agent         string  `gorm:"column:agent;size:100" json:"agent"`
	Brieforg      string  `gorm:"column:brieforg;size:20" json:"brieforg"`
	ClientId      int     `gorm:"column:client_id" json:"client_id"`
	ClientName    string  `gorm:"column:client_name;size:120" json:"client_name"`
	ClientAddress string  `gorm:"column:client_address;size:150" json:"client_address"`
	VidDoc        string  `gorm:"column:vid_doc;size:100" json:"vid_doc"`
	StartAt       string  `gorm:"column:start_at" json:"start_at"`
	FinishAt      string  `gorm:"column:finish_at" json:"finish_at"`
	Done          bool    `gorm:"column:done" json:"done"`
	Status        int     `gorm:"column:status" json:"status"`
	UserID        uint    `gorm:"column:user_id" json:"user_id"`
	EmployeeID    uint    `gorm:"column:employee_id" json:"employee_id"`
	Check         bool    `gorm:"column:check" json:"check"`
}

type AssemblyOrder struct {
	OrderUid        int     `gorm:"column:order_uid;not null;unique" json:"order_uid"`
	OrderDate       string  `gorm:"column:order_date;size:19" json:"order_date"`
	OrderSum        float32 `gorm:"column:order_sum" json:"order_sum"`
	FolioNum        int     `gorm:"column:folio_num" json:"folio_num"`
	FolioDate       string  `gorm:"column:folio_date" json:"folio_date"`
	UnicumNum       int     `gorm:"column:unicum_num" json:"unicum_num"`
	FolioSum        float64 `gorm:"column:folio_sum" json:"folio_sum"`
	UserID          uint    `gorm:"column:user_id" json:"user_id"`
	EmployeeID      uint    `gorm:"column:employee_id" json:"employee_id"`
	CreatedAt       time.Time
	AssemblyDate    string `gorm:"column:assembly_date;size:19" json:"assembly_date"`
	DateDiffMinutes string `gorm:"column:date_diff_minutes;size:19" json:"date_diff_minutes"`
	DateDiffHours   string `gorm:"column:date_diff_hours;size:19" json:"date_diff_hours"`
	UserName        string `gorm:"column:user_name" json:"user_name"`
	EmployeeName    string `gorm:"column:employee_name" json:"employee_name"`
	ClientName      string `gorm:"column:client_name;size:120" json:"client_name"`
	VidDoc          string `gorm:"column:vid_doc;size:100" json:"vid_doc"`
}

func (Order) TableName() string {
	return "orders"
}
