package region_api

import (
	"TPC-EDM-Server/common"
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/models"
	"github.com/gin-gonic/gin"
)

type RegionListRequest struct {
	common.PageInfo
}

type RegionListResponse struct {
	models.RegionModel
}

func (RegionApi) RegionListView(c *gin.Context) {
	var cr RegionListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	list, count, err := common.ListQuery(models.RegionModel{}, common.Options{
		PageInfo: cr.PageInfo,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, count, c)
	return
}
