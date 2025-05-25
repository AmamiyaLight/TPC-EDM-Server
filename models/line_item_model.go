package models

import "time"

type LineItemModel struct {
	OrderKey      uint      `json:"order_key" gorm:"column:L_ORDERKEY;primaryKey"`
	PartKey       uint      `json:"part_key" gorm:"column:L_PARTKEY"`
	SuppKey       uint      `json:"supp_key" gorm:"column:L_SUPPKEY"`
	LineNumber    int       `json:"line_number" gorm:"column:L_LINENUMBER;primaryKey"`
	Quantity      float64   `json:"quantity" gorm:"column:L_QUANTITY;type:decimal(15,2)"`
	ExtendedPrice float64   `json:"extended_price" gorm:"column:L_EXTENDEDPRICE;type:decimal(15,2)"`
	Discount      float64   `json:"discount" gorm:"column:L_DISCOUNT;type:decimal(15,2)"`
	Tax           float64   `json:"tax" gorm:"column:L_TAX;type:decimal(15,2)"`
	ReturnFlag    string    `json:"return_flag" gorm:"column:L_RETURNFLAG;type:char(1)"`
	LineStatus    string    `json:"line_status" gorm:"column:L_LINESTATUS;type:char(1)"`
	ShipDate      time.Time `json:"ship_date" gorm:"column:L_SHIPDATE;type:date"`
	CommitDate    time.Time `json:"commit_date" gorm:"column:L_COMMITDATE;type:date"`
	ReceiptDate   time.Time `json:"receipt_date" gorm:"column:L_RECEIPTDATE;type:date"`
	ShipInstruct  string    `json:"ship_instruct" gorm:"column:L_SHIPINSTRUCT;type:char(25)"`
	ShipMode      string    `json:"ship_mode" gorm:"column:L_SHIPMODE;type:char(10)"`
	Comment       string    `json:"comment" gorm:"column:L_COMMENT;type:varchar(44)"`
}
