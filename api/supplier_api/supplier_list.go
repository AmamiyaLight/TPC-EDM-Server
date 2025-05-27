package supplier_api

import (
	"TPC-EDM-Server/common"
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/models"
	"github.com/gin-gonic/gin"
)

type SupplierListRequest struct {
	common.PageInfo
}

type SupplierListResponse struct {
	models.SupplierModel
}

func (SupplierApi) SupplierListView(c *gin.Context) {
	var cr SupplierListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	list, count, err := common.ListQuery(models.SupplierModel{}, common.Options{
		PageInfo: cr.PageInfo,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, count, c)
	return
}
