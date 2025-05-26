package part_supp_api

import (
	"TPC-H-EDM-Server/common"
	"TPC-H-EDM-Server/common/res"
	"TPC-H-EDM-Server/models"
	"github.com/gin-gonic/gin"
)

type PartSuppListRequest struct {
	common.PageInfo
}

type PartSuppListResponse struct {
	models.PartSuppModel
}

func (PartSuppApi) PartSuppListView(c *gin.Context) {
	var cr PartSuppListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	list, count, err := common.ListQuery(models.PartSuppModel{}, common.Options{
		PageInfo: cr.PageInfo,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, count, c)
	return
}
