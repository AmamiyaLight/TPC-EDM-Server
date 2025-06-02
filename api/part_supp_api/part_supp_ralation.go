package part_supp_api

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"github.com/gin-gonic/gin"
	"time"
)

type PartSuppRelationResponse struct {
	Brand string `json:"brand" gorm:"column:P_BRAND"`
	Size  int    `json:"size" gorm:"column:P_SIZE"`
	Type  string `json:"type" gorm:"column:P_TYPE"`
	Count int64  `json:"count" gorm:"column:supplier_cnt"`
}

// 新增：包含完整响应数据的结构体
type PartSuppQueryResponse struct {
	Results     []PartSuppRelationResponse `json:"results"`      // 查询结果
	Duration    string                     `json:"duration"`     // 查询执行时间
	ExplainPlan string                     `json:"explain_plan"` // 执行计划详情
}

func (PartSuppApi) PartSuppRelationView(c *gin.Context) {
	// 构建查询SQL
	query := `
	SELECT
		p.P_BRAND,
		p.P_TYPE,
		p.P_SIZE,
		COUNT(DISTINCT ps.PS_SUPPKEY) AS supplier_cnt
	FROM part_models p
	JOIN part_supp_models ps ON p.P_PARTKEY = ps.PS_PARTKEY
	LEFT JOIN supplier_models s 
		ON ps.PS_SUPPKEY = s.S_SUPPKEY
		AND s.S_COMMENT LIKE '%Customer%Complaints%'  -- 合并到JOIN条件
	WHERE
		p.P_BRAND <> 'Brand#45'
		AND p.P_TYPE NOT LIKE 'MEDIUM POLISHED%'
		AND p.P_SIZE IN (49, 14, 23, 45, 19, 3, 36, 9)
		AND s.S_SUPPKEY IS NULL  -- 排除匹配项
	GROUP BY p.P_BRAND, p.P_TYPE, p.P_SIZE
	ORDER BY supplier_cnt DESC, p.P_BRAND, p.P_TYPE, p.P_SIZE;
	`

	// 1. 执行原始查询并计时
	startTime := time.Now()
	var results []PartSuppRelationResponse
	err := global.DB.Raw(query).Scan(&results).Error
	duration := time.Since(startTime).String() // 获取耗时

	if err != nil {
		res.FailWithError(err, c)
		return
	}

	// 2. 获取执行计划
	explainQuery := "EXPLAIN FORMAT=JSON " + query
	var explainResult struct {
		Explain string `gorm:"column:EXPLAIN"`
	}
	err = global.DB.Raw(explainQuery).Scan(&explainResult).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	// 3. 构建完整响应
	response := PartSuppQueryResponse{
		Results:     results,
		Duration:    duration,
		ExplainPlan: explainResult.Explain,
	}

	res.OkWithData(response, c)
}
