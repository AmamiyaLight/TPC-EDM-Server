package models

type RegionModel struct {
	RegionKey uint   `json:"region_key" gorm:"column:R_REGIONKEY;primaryKey"`
	Name      string `json:"name" gorm:"column:R_NAME;type:char(25)"`
	Comment   string `json:"comment" gorm:"column:R_COMMENT;type:varchar(152)"`
}
