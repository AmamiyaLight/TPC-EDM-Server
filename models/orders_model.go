package models

import "time"

type OrdersModel struct {
	OrderKey      uint          `json:"order_key" gorm:"column:O_ORDERKEY;primaryKey"`
	CustKey       uint          `json:"cust_key" gorm:"column:O_CUSTKEY"`
	OrderStatus   string        `json:"order_status" gorm:"column:O_ORDERSTATUS;type:char(1)"`
	TotalPrice    float64       `json:"total_price" gorm:"column:O_TOTALPRICE;type:decimal(15,2)"`
	OrderDate     time.Time     `json:"order_date" gorm:"column:O_ORDERDATE;type:date"`
	OrderPriority string        `json:"order_priority" gorm:"column:O_ORDERPRIORITY;type:char(15)"`
	Clerk         string        `json:"clerk" gorm:"column:O_CLERK;type:char(15)"`
	ShipPriority  int           `json:"ship_priority" gorm:"column:O_SHIPPRIORITY"`
	Comment       string        `json:"comment" gorm:"column:O_COMMENT;type:varchar(79)"`
	CustomerModel CustomerModel `json:"-" gorm:"foreignKey:CustKey"`
}
