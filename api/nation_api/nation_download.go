package nation_api

import (
	"TPC-EDM-Server/common/file"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/utils/parse_utils"
	"github.com/gin-gonic/gin"
)

func (api *NationApi) NationDownloadView(c *gin.Context) {
	file.DownloadHandler(
		c,
		"Nation.csv",
		[]string{"N_NATIONKEY", "N_NAME", "N_REGIONKEY", "N_COMMENT"},
		func() int64 {
			var count int64
			global.DB.Model(&models.NationModel{}).Count(&count)
			return count
		},
		func(offset, limit int) ([]models.NationModel, error) {
			var Nations []models.NationModel
			err := global.DB.Order("N_NATIONKEY").Offset(offset).Limit(limit).Find(&Nations).Error
			return Nations, err
		},
		func(Nation models.NationModel) []string {
			return []string{
				parse_utils.StrConvUInt(Nation.NationKey),
				Nation.Name,
				parse_utils.StrConvUInt(Nation.RegionKey),
				Nation.Comment,
			}
		},
	)
}
