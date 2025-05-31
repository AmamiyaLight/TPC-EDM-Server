package supplier_api

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"github.com/gin-gonic/gin"
)

// 请求参数结

// 响应数据结构体
type SupplierRevenue struct {
	SuppKey  int64   `json:"supp_key" gorm:"column:s_suppkey"`
	Name     string  `json:"name" gorm:"column:s_name"`
	Address  string  `json:"address" gorm:"column:s_address"`
	Phone    string  `json:"phone" gorm:"column:s_phone"`
	TotalRev float64 `json:"total_rev" gorm:"column:total_revenue"`
}

func (SupplierApi) TopSupplierView(c *gin.Context) {

	query := `
	WITH revenue0 AS (
    SELECT
        l_suppkey AS supplier_no,
        SUM(l_extendedprice * (1 - l_discount)) AS total_revenue,
        MAX(SUM(l_extendedprice * (1 - l_discount))) OVER () AS max_revenue
    FROM line_item_models
    WHERE l_shipdate >= ?
      AND l_shipdate < ?
    GROUP BY l_suppkey
	)
	SELECT
		s.s_suppkey,
		s.s_name,
		s.s_address,
		s.s_phone,
		r.total_revenue
	FROM supplier_models s
			 JOIN revenue0 r
				  ON s.s_suppkey = r.supplier_no
	WHERE r.total_revenue = r.max_revenue
	ORDER BY s.s_suppkey;
	`

	// 执行查询
	var results SupplierRevenue
	err := global.DB.Raw(
		query,
		"1997-01-28",
		"1997-04-28",
	).Scan(&results).Error

	if err != nil {
		res.FailWithError(err, c)
		return
	}

	res.OkWithData(results, c)
}
