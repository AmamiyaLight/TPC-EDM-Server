package part_supp_api

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"github.com/gin-gonic/gin"
)

type PartSuppRelationResponse struct {
	Brand string `json:"brand" gorm:"column:P_BRAND"`
	Size  int    `json:"size" gorm:"column:P_SIZE"`
	Type  string `json:"type" gorm:"column:P_TYPE"`
	Count int64  `json:"count" gorm:"column:supplier_cnt"`
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

	// 执行查询（返回多条记录）
	var results []PartSuppRelationResponse
	err := global.DB.Raw(query).Scan(&results).Error

	if err != nil {
		res.FailWithError(err, c)
		return
	}

	res.OkWithData(results, c)
}
