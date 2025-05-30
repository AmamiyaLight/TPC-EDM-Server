package supplier_api

import (
	"TPC-EDM-Server/common/file"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/utils/parse_utils"
	"github.com/gin-gonic/gin"
)

func (api *SupplierApi) SupplierDownloadView(c *gin.Context) {
	file.DownloadHandler(
		c,
		"Supplier.csv",
		[]string{"S_SUPPKEY",
			"S_NAME",
			"S_ADDRESS",
			"S_NATIONKEY",
			"S_PHONE",
			"S_ACCTBAL",
			"S_COMMENT"},
		func() int64 {
			var count int64
			global.DB.Model(&models.SupplierModel{}).Count(&count)
			return count
		},
		func(offset, limit int) ([]models.SupplierModel, error) {
			var Suppliers []models.SupplierModel
			err := global.DB.Order("S_SUPPKEY").Offset(offset).Limit(limit).Find(&Suppliers).Error
			return Suppliers, err
		},
		func(Supplier models.SupplierModel) []string {
			return []string{
				parse_utils.StrConvUInt(Supplier.SuppKey),
				Supplier.Name,
				Supplier.Address,
				parse_utils.StrConvUInt(Supplier.NationKey),
				Supplier.Phone,
				parse_utils.StrConvFloat(Supplier.AcctBal),
				Supplier.Comment,
			}
		},
	)
}
