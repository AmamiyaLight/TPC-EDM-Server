package models

type PartModel struct {
	PartKey     uint    `json:"part_key" gorm:"column:P_PARTKEY;primaryKey"`
	Name        string  `json:"name" gorm:"column:P_NAME;type:varchar(55)"`
	Mfgr        string  `json:"mfgr" gorm:"column:P_MFGR;type:char(25)"`
	Brand       string  `json:"brand" gorm:"column:P_BRAND;type:char(10)"`
	Type        string  `json:"type" gorm:"column:P_TYPE;type:varchar(25)"`
	Size        int     `json:"size" gorm:"column:P_SIZE"`
	Container   string  `json:"container" gorm:"column:P_CONTAINER;type:char(10)"`
	RetailPrice float64 `json:"retail_price" gorm:"column:P_RETAILPRICE;type:decimal(15,2)"`
	Comment     string  `json:"comment" gorm:"column:P_COMMENT;type:varchar(23)"`
}
