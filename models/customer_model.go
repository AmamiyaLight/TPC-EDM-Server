package models

type CustomerModel struct {
	CustKey    uint    `json:"cust_key" gorm:"column:C_CUSTKEY;primaryKey"`
	Name       string  `json:"name" gorm:"column:C_NAME;type:varchar(25)"`
	Address    string  `json:"address" gorm:"column:C_ADDRESS;type:varchar(40)"`
	NationKey  uint    `json:"nation_key" gorm:"column:C_NATIONKEY"`
	Phone      string  `json:"phone" gorm:"column:C_PHONE;type:char(15)"`
	AcctBal    float64 `json:"acct_bal" gorm:"column:C_ACCTBAL;type:decimal(15,2)"`
	MktSegment string  `json:"mkt_segment" gorm:"column:C_MKTSEGMENT;type:char(10)"`
	Comment    string  `json:"comment" gorm:"column:C_COMMENT;type:varchar(117)"`
}
