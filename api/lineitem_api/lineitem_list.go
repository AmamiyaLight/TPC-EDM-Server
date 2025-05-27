package lineitem_api

import (
	"TPC-EDM-Server/common"
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/models"
	"github.com/gin-gonic/gin"
)

type LineItemListRequest struct {
	common.PageInfo
}

type LineItemListResponse struct {
	models.LineItemModel
}

func (LineItemApi) LineItemListView(c *gin.Context) {
	var cr LineItemListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	list, count, err := common.ListQuery(models.LineItemModel{}, common.Options{
		PageInfo: cr.PageInfo,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, count, c)
	return
}
