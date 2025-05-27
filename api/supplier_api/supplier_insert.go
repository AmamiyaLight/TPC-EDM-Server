package supplier_api

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

func (SupplierApi) SupplierInsertView(c *gin.Context) {
	startTime := time.Now()

	total, err := file.ProcessFileInsert(c, parseSupplierLine, 1000)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg(fmt.Sprintf("导入成功,共导入%d条数据,耗时%s",
		total, time.Since(startTime).String()), c)
}
func parseSupplierLine(line string) (models.SupplierModel, error) {
	fields := strings.Split(line, "|")
	if len(fields) < 7 {
		return models.SupplierModel{}, errors.New("字段不足")
	}
	return models.SupplierModel{
		SuppKey:   parse_utils.ParseUint(fields[0]),
		Name:      fields[1],
		Address:   fields[2],
		NationKey: parse_utils.ParseUint(fields[3]),
		Phone:     fields[4],
		AcctBal:   parse_utils.ParseFloat64(fields[5]),
		Comment:   fields[6],
	}, nil
}
