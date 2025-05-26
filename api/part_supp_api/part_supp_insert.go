package part_supp_api

import (
	"TPC-H-EDM-Server/common/file"
	"TPC-H-EDM-Server/common/res"
	"TPC-H-EDM-Server/models"
	"TPC-H-EDM-Server/utils/parse_utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func (PartSuppApi) PartSuppInsertView(c *gin.Context) {
	startTime := time.Now()

	total, err := file.ProcessFileInsert(c, parsePartSuppLine, 1000)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg(fmt.Sprintf("导入成功,共导入%d条数据,耗时%s",
		total, time.Since(startTime).String()), c)
}
func parsePartSuppLine(line string) (models.PartSuppModel, error) {
	fields := strings.Split(line, "|")
	if len(fields) < 5 {
		return models.PartSuppModel{}, errors.New("字段不足")
	}
	return models.PartSuppModel{
		PartKey:    parse_utils.ParseUint(fields[0]),
		SuppKey:    parse_utils.ParseUint(fields[1]),
		AvailQty:   parse_utils.ParseIntUtil(fields[2]),
		SupplyCost: parse_utils.ParseFloat64(fields[3]),
		Comment:    fields[4],
	}, nil
}
