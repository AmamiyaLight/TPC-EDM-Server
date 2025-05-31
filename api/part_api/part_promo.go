package part_api

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"github.com/gin-gonic/gin"
	"time"
)

type PartPromoRequest struct {
	Date string `form:"date" binding:"required"`
}

func (PartApi) PartPromoView(c *gin.Context) {
	var cr PartPromoRequest
	if err := c.ShouldBind(&cr); err != nil {
		res.FailWithError(err, c)
		return
	}

	startDate, _ := time.Parse("2006-01-02", cr.Date)
	endDate := startDate.AddDate(0, 1, 0)

	// 使用子查询限制数据集大小
	query := `
    SELECT COALESCE(
        100.00 * 
        SUM(CASE WHEN p.p_type LIKE 'PROMO%%' 
                 THEN l.l_extendedprice * (1 - l.l_discount) 
                 ELSE 0 END) /
        NULLIF(SUM(l.l_extendedprice * (1 - l.l_discount)), 0),
    0) AS percent
    FROM line_item_models l
    JOIN part_models p ON l.l_partkey = p.p_partkey
    WHERE l.l_shipdate >= ?
      AND l.l_shipdate < ?
	`

	var promoPercent float64
	err := global.DB.Raw(query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Scan(&promoPercent).Error

	if err != nil {
		res.FailWithError(err, c)
		return
	}

	res.OkWithData(promoPercent, c)
}
