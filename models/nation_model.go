package models

type NationModel struct {
	NationKey uint   `json:"nation_key" gorm:"column:N_NATIONKEY;primaryKey"`
	Name      string `json:"name" gorm:"column:N_NAME;type:char(25)"`
	RegionKey uint   `json:"region_key" gorm:"column:N_REGIONKEY"`
	Comment   string `json:"comment" gorm:"column:N_COMMENT;type:varchar(152)"`
}
