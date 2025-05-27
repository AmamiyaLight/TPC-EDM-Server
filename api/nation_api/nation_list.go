package nation_api

import (
	"TPC-EDM-Server/common"
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/models"
	"github.com/gin-gonic/gin"
)

type NationListRequest struct {
	common.PageInfo
}

type NationListResponse struct {
	models.NationModel
}

func (NationApi) NationListView(c *gin.Context) {
	var cr NationListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	list, count, err := common.ListQuery(models.NationModel{}, common.Options{
		PageInfo: cr.PageInfo,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, count, c)
	return
}
