package part_api

import (
	"TPC-EDM-Server/common/file"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/utils/parse_utils"
	"github.com/gin-gonic/gin"
)

func (api *PartApi) PartDownloadView(c *gin.Context) {
	file.DownloadHandler(
		c,
		"Part.csv",
		[]string{"P_PARTKEY",
			"P_NAME",
			"P_MFGR",
			"P_BRAND",
			"P_TYPE",
			"P_SIZE",
			"P_CONTAINER",
			"P_RETAILPRICE",
			"P_COMMENT"},
		func() int64 {
			var count int64
			global.DB.Model(&models.PartModel{}).Count(&count)
			return count
		},
		func(offset, limit int) ([]models.PartModel, error) {
			var Parts []models.PartModel
			err := global.DB.Order("P_PARTKEY").Offset(offset).Limit(limit).Find(&Parts).Error
			return Parts, err
		},
		func(Part models.PartModel) []string {
			return []string{
				parse_utils.StrConvUInt(Part.PartKey),
				Part.Name,
				Part.Mfgr,
				Part.Brand,
				Part.Type,
				parse_utils.StrConvInt(Part.Size),
				Part.Container,
				parse_utils.StrConvFloat(Part.RetailPrice),
				Part.Comment,
			}
		},
	)
}
