package models

type SupplierModel struct {
	SuppKey   uint    `json:"supp_key" gorm:"column:S_SUPPKEY;primaryKey"`
	Name      string  `json:"name" gorm:"column:S_NAME;type:char(25)"`
	Address   string  `json:"address" gorm:"column:S_ADDRESS;type:varchar(40)"`
	NationKey uint    `json:"nation_key" gorm:"column:S_NATIONKEY"`
	Phone     string  `json:"phone" gorm:"column:S_PHONE;type:char(15)"`
	AcctBal   float64 `json:"acct_bal" gorm:"column:S_ACCTBAL;type:decimal(15,2)"`
	Comment   string  `json:"comment" gorm:"column:S_COMMENT;type:varchar(101)"`
}
