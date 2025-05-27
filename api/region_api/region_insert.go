package region_api

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

func (RegionApi) RegionInsertView(c *gin.Context) {
	startTime := time.Now()

	total, err := file.ProcessFileInsert(c, parseRegionLine, 1000)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg(fmt.Sprintf("导入成功,共导入%d条数据,耗时%s",
		total, time.Since(startTime).String()), c)
}
func parseRegionLine(line string) (models.RegionModel, error) {
	fields := strings.Split(line, "|")
	if len(fields) < 3 {
		return models.RegionModel{}, errors.New("字段不足")
	}
	return models.RegionModel{
		RegionKey: parse_utils.ParseUint(fields[0]),
		Name:      fields[1],
		Comment:   fields[2],
	}, nil
}
