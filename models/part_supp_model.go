package models

type PartSuppModel struct {
	PartKey    uint    `json:"part_key" gorm:"column:PS_PARTKEY;primaryKey"`
	SuppKey    uint    `json:"supp_key" gorm:"column:PS_SUPPKEY;primaryKey"`
	AvailQty   int     `json:"avail_qty" gorm:"column:PS_AVAILQTY"`
	SupplyCost float64 `json:"supply_cost" gorm:"column:PS_SUPPLYCOST;type:decimal(15,2)"`
	Comment    string  `json:"comment" gorm:"column:PS_COMMENT;type:varchar(199)"`
}
