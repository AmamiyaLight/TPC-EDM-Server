package part_api

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

func (PartApi) PartInsertView(c *gin.Context) {
	startTime := time.Now()

	total, err := file.ProcessFileInsert(c, parsePartLine, 1000)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg(fmt.Sprintf("导入成功,共导入%d条数据,耗时%s",
		total, time.Since(startTime).String()), c)
}
func parsePartLine(line string) (models.PartModel, error) {
	fields := strings.Split(line, "|")
	if len(fields) < 5 {
		return models.PartModel{}, errors.New("字段不足")
	}
	return models.PartModel{
		PartKey:     parse_utils.ParseUint(fields[0]),
		Name:        fields[1],
		Mfgr:        fields[2],
		Brand:       fields[3],
		Type:        fields[4],
		Size:        parse_utils.ParseIntUtil(fields[5]),
		Container:   fields[6],
		RetailPrice: parse_utils.ParseFloat64(fields[7]),
		Comment:     fields[8],
	}, nil
}
