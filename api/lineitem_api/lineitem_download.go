package lineitem_api

import (
	"TPC-EDM-Server/common/file"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/utils/parse_utils"
	"github.com/gin-gonic/gin"
)

func (api *LineItemApi) LineItemDownloadView(c *gin.Context) {
	file.DownloadHandler(
		c,
		"LineItem.csv",
		[]string{"L_ORDERKEY",
			"L_PARTKEY",
			"L_SUPPKEY",
			"L_LINENUMBER",
			"L_QUANTITY",
			"L_EXTENDEDPRICE",
			"L_DISCOUNT",
			"L_TAX",
			"L_RETURNFLAG",
			"L_LINESTATUS",
			"L_SHIPDATE",
			"L_COMMITDATE",
			"L_RECEIPTDATE",
			"L_SHIPINSTRUCT",
			"L_SHIPMODE",
			"L_COMMENT"},
		func() int64 {
			var count int64
			global.DB.Model(&models.LineItemModel{}).Count(&count)
			return count
		},
		func(offset, limit int) ([]models.LineItemModel, error) {
			var LineItems []models.LineItemModel
			err := global.DB.Order("L_ORDERKEY ASC, L_LINENUMBER ASC").Offset(offset).Limit(limit).Find(&LineItems).Error
			return LineItems, err
		},
		func(LineItem models.LineItemModel) []string {
			return []string{
				parse_utils.StrConvUInt(LineItem.OrderKey),
				parse_utils.StrConvUInt(LineItem.PartKey),
				parse_utils.StrConvUInt(LineItem.SuppKey),
				parse_utils.StrConvInt(LineItem.LineNumber),
				parse_utils.StrConvFloat(LineItem.Quantity),
				parse_utils.StrConvFloat(LineItem.ExtendedPrice),
				parse_utils.StrConvFloat(LineItem.Discount),
				parse_utils.StrConvFloat(LineItem.Tax),
				LineItem.ReturnFlag,
				LineItem.LineStatus,
				parse_utils.StrConvTime(LineItem.ShipDate),
				parse_utils.StrConvTime(LineItem.CommitDate),
				parse_utils.StrConvTime(LineItem.ReceiptDate),
				LineItem.ShipInstruct,
				LineItem.ShipMode,
				LineItem.Comment,
			}
		},
	)
}
