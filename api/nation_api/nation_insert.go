package nation_api

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

func (NationApi) NationInsertView(c *gin.Context) {
	startTime := time.Now()

	total, err := file.ProcessFileInsert(c, parseNationLine, 1000)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg(fmt.Sprintf("导入成功,共导入%d条数据,耗时%s",
		total, time.Since(startTime).String()), c)
}
func parseNationLine(line string) (models.NationModel, error) {
	fields := strings.Split(line, "|")
	if len(fields) < 4 {
		return models.NationModel{}, errors.New("字段不足")
	}
	return models.NationModel{
		NationKey: parse_utils.ParseUint(fields[0]),
		Name:      fields[1],
		RegionKey: parse_utils.ParseUint(fields[2]),
		Comment:   fields[3],
	}, nil
}
