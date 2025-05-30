package region_api

import (
	"TPC-EDM-Server/common/file"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/utils/parse_utils"
	"github.com/gin-gonic/gin"
)

func (api *RegionApi) RegionDownloadView(c *gin.Context) {
	file.DownloadHandler(
		c,
		"Region.csv",
		[]string{"R_REGIONKEY", "R_NAME", "R_COMMENT"},
		func() int64 {
			var count int64
			global.DB.Model(&models.RegionModel{}).Count(&count)
			return count
		},
		func(offset, limit int) ([]models.RegionModel, error) {
			var Regions []models.RegionModel
			err := global.DB.Order("R_REGIONKEYnation_download.go").Offset(offset).Limit(limit).Find(&Regions).Error
			return Regions, err
		},
		func(Region models.RegionModel) []string {
			return []string{
				parse_utils.StrConvUInt(Region.RegionKey),
				Region.Name,
				Region.Comment,
			}
		},
	)
}
