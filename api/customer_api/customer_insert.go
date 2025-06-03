package customer_api

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

func (CustomerApi) CustomerInsertView(c *gin.Context) {
	startTime := time.Now()

	total, err := file.ProcessFileInsert(c, parseCustomerLine, 1000)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg(fmt.Sprintf("导入成功,共导入%d条数据,耗时%s",
		total, time.Since(startTime).String()), c)
}
func parseCustomerLine(line string) (models.CustomerModel, error) {
	fields := strings.Split(line, "|")
	if len(fields) < 5 {
		return models.CustomerModel{}, errors.New("字段不足")
	}
	item := models.CustomerModel{
		CustKey:    parse_utils.ParseUint(fields[0]),
		Name:       fields[1],
		Address:    fields[2],
		NationKey:  parse_utils.ParseUint(fields[3]),
		Phone:      fields[4],
		AcctBal:    parse_utils.ParseFloat64(fields[5]),
		MktSegment: fields[6],
		Comment:    fields[7],
	}
	if item.CustKey == 0 || item.NationKey == 0 {
		return models.CustomerModel{}, errors.New("主键或外键校验失败")
	}
	if item.AcctBal < 0 {
		return models.CustomerModel{}, errors.New("账户余额不能为负")
	}

	return item, nil
}
