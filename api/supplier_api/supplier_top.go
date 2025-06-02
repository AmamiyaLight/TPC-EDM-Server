package supplier_api

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

type SupplierRevenue struct {
	SuppKey  int64   `json:"supp_key" gorm:"column:s_suppkey"`
	Name     string  `json:"name" gorm:"column:s_name"`
	Address  string  `json:"address" gorm:"column:s_address"`
	Phone    string  `json:"phone" gorm:"column:s_phone"`
	TotalRev float64 `json:"total_rev" gorm:"column:total_revenue"`
}

// 新增：包含完整响应数据的结构体
type QueryResponse struct {
	Results     []SupplierRevenue `json:"results"`      // 查询结果
	Duration    string            `json:"duration"`     // 查询执行时间
	ExplainPlan string            `json:"explain_plan"` // 执行计划详情
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

	params := []interface{}{"1997-01-28", "1997-04-28"}

	startTime := time.Now()
	var results []SupplierRevenue
	err := global.DB.Raw(query, params...).Scan(&results).Error
	duration := time.Since(startTime).String() // 获取耗时

	if err != nil {
		res.FailWithError(err, c)
		return
	}

	explainQuery := "EXPLAIN " + query
	var explainRows []map[string]interface{}
	err = global.DB.Raw(explainQuery, params...).Scan(&explainRows).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	explainJSON, _ := json.MarshalIndent(explainRows, "", "  ")

	response := QueryResponse{
		Results:     results,
		Duration:    duration,
		ExplainPlan: string(explainJSON),
	}

	res.OkWithData(response, c)
}
