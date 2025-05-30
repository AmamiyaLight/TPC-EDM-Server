package part_supp_api

import (
	"TPC-EDM-Server/common/file"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/utils/parse_utils"
	"github.com/gin-gonic/gin"
)

func (api *PartSuppApi) PartSuppDownloadView(c *gin.Context) {
	file.DownloadHandler(
		c,
		"PartSupp.csv",
		[]string{"PS_PARTKEY",
			"PS_SUPPKEY",
			"PS_AVAILQTY",
			"PS_SUPPLYCOST",
			"PS_COMMENT"},
		func() int64 {
			var count int64
			global.DB.Model(&models.PartSuppModel{}).Count(&count)
			return count
		},
		func(offset, limit int) ([]models.PartSuppModel, error) {
			var PartSupps []models.PartSuppModel
			err := global.DB.Order("PS_PARTKEY ASC, PS_SUPPKEY ASC").Offset(offset).Limit(limit).Find(&PartSupps).Error
			return PartSupps, err
		},
		func(PartSupp models.PartSuppModel) []string {
			return []string{
				parse_utils.StrConvUInt(PartSupp.PartKey),
				parse_utils.StrConvUInt(PartSupp.SuppKey),
				parse_utils.StrConvInt(PartSupp.AvailQty),
				parse_utils.StrConvFloat(PartSupp.SupplyCost),
				PartSupp.Comment,
			}
		},
	)
}
