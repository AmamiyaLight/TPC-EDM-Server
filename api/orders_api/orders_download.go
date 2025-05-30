package orders_api

import (
	"TPC-EDM-Server/common/file"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/utils/parse_utils"
	"github.com/gin-gonic/gin"
)

func (api *OrdersApi) OrdersDownloadView(c *gin.Context) {
	file.DownloadHandler(
		c,
		"Orders.csv",
		[]string{"O_ORDERKEY",
			"O_CUSTKEY",
			"O_ORDERSTATUS",
			"O_TOTALPRICE",
			"O_ORDERDATE",
			"O_ORDERPRIORITY",
			"O_CLERK",
			"O_SHIPPRIORITY",
			"O_COMMENT"},
		func() int64 {
			var count int64
			global.DB.Model(&models.OrdersModel{}).Count(&count)
			return count
		},
		func(offset, limit int) ([]models.OrdersModel, error) {
			var Orderss []models.OrdersModel
			err := global.DB.Order("O_ORDERKEY").Offset(offset).Limit(limit).Find(&Orderss).Error
			return Orderss, err
		},
		func(Orders models.OrdersModel) []string {
			return []string{
				parse_utils.StrConvUInt(Orders.OrderKey),
				parse_utils.StrConvUInt(Orders.CustKey),
				Orders.OrderStatus,
				parse_utils.StrConvFloat(Orders.TotalPrice),
				parse_utils.StrConvTime(Orders.OrderDate),
				Orders.OrderPriority,
				Orders.Clerk,
				parse_utils.StrConvInt(Orders.ShipPriority),
				Orders.Comment,
			}
		},
	)
}
