package customer_api

import (
	"TPC-EDM-Server/common"
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"fmt"
	"github.com/gin-gonic/gin"
)

type CustomerListRequest struct {
	common.PageInfo
	Name   string `form:"name"`
	Nation string `form:"nation"`
}

type CustomerListResponse struct {
	models.CustomerModel
}

func (CustomerApi) CustomerListView(c *gin.Context) {
	var cr CustomerListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	if cr.Name != "" || cr.Nation != "" {
		joins := []string{}
		query := global.DB.Where("")
		if cr.Name != "" {
			query = query.Where("customer_models.C_NAME LIKE ?", fmt.Sprintf("%%%s%%", cr.Name))
		}
		if cr.Nation != "" {
			joins = []string{
				"LEFT JOIN nation_models ON customer_models.C_NATIONKEY = nation_models.N_NATIONKEY",
			}
			query = query.Where("nation_models.N_NAME LIKE ?", fmt.Sprintf("%%%s%%", cr.Nation))
		}
		list, count, err := common.ListQuery(models.CustomerModel{}, common.Options{
			PageInfo:      cr.PageInfo,
			Joins:         joins,
			AccuracyCount: true,
			Where:         query,
		})
		if err != nil {
			res.FailWithMsg(fmt.Sprintf("查询错误%s", err), c)
			return
		}
		res.OkWithList(list, count, c)
		return
	}

	list, count, err := common.ListQuery(models.CustomerModel{}, common.Options{
		PageInfo: cr.PageInfo,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, count, c)
	return
}
