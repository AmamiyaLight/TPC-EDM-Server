package customer_api

import (
	"TPC-EDM-Server/common/file"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (api *CustomerApi) CustomerDownloadView(c *gin.Context) {
	file.DownloadHandler(
		c,
		"customer.csv",
		[]string{"C_CUSTKEY", "C_NAME", "C_ADDRESS", "C_NATIONKEY", "C_PHONE", "C_ACCTBAL", "C_MKTSEGMENT", "C_COMMENT"},
		func() int64 {
			var count int64
			global.DB.Model(&models.CustomerModel{}).Count(&count)
			return count
		},
		func(offset, limit int) ([]models.CustomerModel, error) {
			var customers []models.CustomerModel
			err := global.DB.Order("C_CUSTKEY").Offset(offset).Limit(limit).Find(&customers).Error
			return customers, err
		},
		func(customer models.CustomerModel) []string {
			return []string{
				strconv.FormatUint(uint64(customer.CustKey), 10),
				customer.Name,
				customer.Address,
				strconv.FormatUint(uint64(customer.NationKey), 10),
				customer.Phone,
				strconv.FormatFloat(customer.AcctBal, 'f', 2, 64),
				customer.MktSegment,
				customer.Comment,
			}
		},
	)
}
