package part_api

import (
	"TPC-EDM-Server/common"
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/models"
	"github.com/gin-gonic/gin"
)

type PartListRequest struct {
	common.PageInfo
}

type PartListResponse struct {
	models.PartModel
}

func (PartApi) PartListView(c *gin.Context) {
	var cr PartListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	list, count, err := common.ListQuery(models.PartModel{}, common.Options{
		PageInfo: cr.PageInfo,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, count, c)
	return
}
