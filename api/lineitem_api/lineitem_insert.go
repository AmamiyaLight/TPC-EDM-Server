package lineitem_api

import (
	"TPC-EDM-Server/common/file"
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/utils/parse_utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func (LineItemApi) LineItemInsertView(c *gin.Context) {
	startTime := time.Now()
	total, err := file.ProcessFileInsert(c, parseLineItemLine, 500)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg(fmt.Sprintf("导入成功,共导入%d条数据,耗时%s",
		total, time.Since(startTime).String()), c)
}

func parseLineItemLine(line string) (models.LineItemModel, error) {
	fields := strings.Split(line, "|")
	if len(fields) < 16 {
		return models.LineItemModel{}, errors.New("字段不足")
	}
	item := models.LineItemModel{
		OrderKey:      parse_utils.ParseUint(fields[0]),
		PartKey:       parse_utils.ParseUint(fields[1]),
		SuppKey:       parse_utils.ParseUint(fields[2]),
		LineNumber:    parse_utils.ParseIntUtil(fields[3]),
		Quantity:      parse_utils.ParseFloat64(fields[4]),
		ExtendedPrice: parse_utils.ParseFloat64(fields[5]),
		Discount:      parse_utils.ParseFloat64(fields[6]),
		Tax:           parse_utils.ParseFloat64(fields[7]),
		ReturnFlag:    fields[8],
		LineStatus:    fields[9],
		ShipDate:      parse_utils.ParseTimeUtil(fields[10]),
		CommitDate:    parse_utils.ParseTimeUtil(fields[11]),
		ReceiptDate:   parse_utils.ParseTimeUtil(fields[12]),
		ShipInstruct:  fields[13],
		ShipMode:      fields[14],
		Comment:       fields[15],
	}
	if item.OrderKey == 0 || item.PartKey == 0 {
		return models.LineItemModel{}, errors.New("主键或外键校验失败")
	}
	if item.LineNumber <= 0 {
		return models.LineItemModel{}, errors.New("行号不能为负")
	}
	if item.Quantity < 0 || item.Quantity > 100 {
		return models.LineItemModel{}, errors.New("quantity字段超限")
	}
	if item.Discount < 0 || item.Tax < 0 || item.ExtendedPrice < 0 {
		return models.LineItemModel{}, errors.New("价格校验失败")
	}
	return item, nil
}
