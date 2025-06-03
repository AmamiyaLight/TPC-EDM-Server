package orders_api

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

func (OrdersApi) OrderInsertView(c *gin.Context) {
	startTime := time.Now()

	total, err := file.ProcessFileInsert(c, parseOrderLine, 1000)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	res.OkWithMsg(fmt.Sprintf("导入成功,共导入%d条数据,耗时%s",
		total, time.Since(startTime).String()), c)
}

func parseOrderLine(line string) (models.OrdersModel, error) {
	fields := strings.Split(line, "|")
	if len(fields) < 9 {
		return models.OrdersModel{}, errors.New("字段不足")
	}
	item := models.OrdersModel{
		OrderKey:      parse_utils.ParseUint(fields[0]),
		CustKey:       parse_utils.ParseUint(fields[1]),
		OrderStatus:   fields[2],
		TotalPrice:    parse_utils.ParseFloat64(fields[3]),
		OrderDate:     parse_utils.ParseTimeUtil(fields[4]),
		OrderPriority: fields[5],
		Clerk:         fields[6],
		ShipPriority:  parse_utils.ParseIntUtil(fields[7]),
		Comment:       fields[8],
	}
	if item.OrderKey == 0 || item.CustKey == 0 {
		return models.OrdersModel{}, errors.New("主键或外键校验失败")
	}
	if item.TotalPrice < 0 {
		return models.OrdersModel{}, errors.New("总价不能为负")
	}
	return item, nil
}
