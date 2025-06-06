package tpcc_api

import (
	"TPC-EDM-Server/global"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewOrderRequest struct {
	W_ID     int `json:"w_id" binding:"required"`
	D_W_ID   int `json:"d_w_id" binding:"required"`
	D_ID     int `json:"d_id" binding:"required"`
	C_W_ID   int `json:"c_w_id" binding:"required"`
	C_D_ID   int `json:"c_d_id" binding:"required"`
	C_ID     int `json:"c_id" binding:"required"`
	O_OL_CNT int `json:"ol_cnt" binding:"required"`
	Items    []struct {
		OL_I_ID        int `json:"ol_i_id"`
		OL_SUPPLY_W_ID int `json:"ol_supply_w_id"`
		OL_QUANTITY    int `json:"ol_quantity"`
	} `json:"items" binding:"required,gt=0,dive"`
}

type StepResult struct {
	StepName    string  `json:"step_name"`
	SQL         string  `json:"sql"`
	TimeMs      float64 `json:"time_ms"`
	ExplainPlan string  `json:"explain_plan"`
	IndexUsed   string  `json:"index_used"`
}

type NewOrderResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Steps   []StepResult `json:"steps"`
}

func (t TpccApi) NewOrderView(c *gin.Context) {
	var req NewOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(req.Items) != req.O_OL_CNT {
		c.JSON(400, gin.H{"error": "OL_CNT doesn't match items count"})
		return
	}

	response := NewOrderResponse{
		Steps: make([]StepResult, 0, 10+req.O_OL_CNT*5),
	}
	var oID int

	// 开始事务
	tx := global.DB.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			response.Success = false
			response.Message = fmt.Sprintf("Panic: %v", r)
		}
	}()

	// 步骤1: 获取客户和仓库信息
	step1 := StepResult{StepName: "get_customer_warehouse"}
	start := time.Now()

	// 创建接收结果的结构体
	type CustomerWarehouseResult struct {
		CDiscount float64
		CLast     string
		CCredit   string
		WTax      float64
	}

	query := tx.Table("customer AS c").
		Select("c.c_discount, c.c_last, c.c_credit, w.w_tax").
		Joins("JOIN warehouse AS w ON w.w_id = c.c_w_id").
		Where("c.c_w_id = ? AND c.c_d_id = ? AND c.c_id = ?", req.C_W_ID, req.C_D_ID, req.C_ID).
		Where("w.w_id = ?", req.W_ID)
	stmt := query.Session(&gorm.Session{DryRun: true}).Statement
	step1.SQL = replacePlaceholders(stmt.SQL.String(), stmt.Vars)

	step1.TimeMs = time.Since(start).Seconds() * 1000
	var result CustomerWarehouseResult
	err := query.Scan(&result).Error
	step1.TimeMs = time.Since(start).Seconds() * 1000
	if recordStep(tx, &step1, err, &response) {
		return
	}

	// 将结果赋值给原有变量
	customer := struct {
		CDiscount float64
		CLast     string
		CCredit   string
	}{
		CDiscount: result.CDiscount,
		CLast:     result.CLast,
		CCredit:   result.CCredit,
	}
	warehouse := struct {
		WTax float64
	}{
		WTax: result.WTax,
	}

	// 步骤2: 获取地区信息
	step2 := StepResult{StepName: "get_district"}
	start = time.Now()
	query = tx.Table("district").
		Select("d_next_o_id, d_tax").
		Where("d_id = ? AND d_w_id = ?", req.D_ID, req.D_W_ID)
	stmt = query.Session(&gorm.Session{DryRun: true}).Statement
	step2.SQL = replacePlaceholders(stmt.SQL.String(), stmt.Vars)

	// 2. 执行查询
	var district struct {
		DNextOID int
		DTax     float64
	}
	err = query.Scan(&district).Error
	step2.TimeMs = time.Since(start).Seconds() * 1000
	if recordStep(tx, &step2, err, &response) {
		return
	}
	// 步骤3: 更新地区订单ID
	step3 := StepResult{StepName: "update_district"}
	start = time.Now()
	sqlStr := "UPDATE district SET d_next_o_id = ? WHERE d_id = ? AND d_w_id = ?"
	args := []interface{}{district.DNextOID + 1, req.D_ID, req.D_W_ID}
	step3.SQL = replacePlaceholders(sqlStr, args)

	// 2. 执行更新
	err = tx.Exec(sqlStr, args...).Error
	step3.TimeMs = time.Since(start).Seconds() * 1000
	if recordStep(tx, &step3, err, &response) {
		return
	}
	oID = district.DNextOID

	// 步骤4: 插入订单
	step4 := StepResult{StepName: "insert_order"}
	start = time.Now()
	oAllLocal := 1
	for _, item := range req.Items {
		if item.OL_SUPPLY_W_ID != req.W_ID {
			oAllLocal = 0
			break
		}
	}

	sqlStr4 := `
		INSERT INTO orders (o_id, o_d_id, o_w_id, o_c_id, o_entry_d, o_ol_cnt, o_all_local)
		VALUES (?, ?, ?, ?, NOW(), ?, ?)`
	args4 := []interface{}{oID, req.D_ID, req.W_ID, req.C_ID, req.O_OL_CNT, oAllLocal}
	step4.SQL = replacePlaceholders(sqlStr4, args4)

	// 2. 执行插入
	err = tx.Exec(sqlStr4, args4...).Error
	step4.TimeMs = time.Since(start).Seconds() * 1000
	if recordStep(tx, &step4, err, &response) {
		return
	}

	// 步骤5: 插入新订单
	step5 := StepResult{StepName: "insert_new_order"}
	start = time.Now()
	sqlStr5 := `
		INSERT INTO new_order (no_o_id, no_d_id, no_w_id)
		VALUES (?, ?, ?)`
	args5 := []interface{}{oID, req.D_ID, req.W_ID}
	step5.SQL = replacePlaceholders(sqlStr5, args5)

	// 2. 执行插入
	err = tx.Exec(sqlStr5, args5...).Error
	step5.TimeMs = time.Since(start).Seconds() * 1000
	if recordStep(tx, &step5, err, &response) {
		return
	}

	// 处理每个订单行
	for olNum, item := range req.Items {
		olNumber := olNum + 1

		// 步骤6: 获取商品信息
		step6 := StepResult{StepName: fmt.Sprintf("get_item_%d", olNumber)}
		start = time.Now()

		// 1. 生成SQL
		query := tx.Table("item").
			Select("i_price, i_name, i_data").
			Where("i_id = ?", item.OL_I_ID)
		stmt := query.Session(&gorm.Session{DryRun: true}).Statement
		step6.SQL = replacePlaceholders(stmt.SQL.String(), stmt.Vars)

		// 2. 执行查询
		var itemInfo struct {
			IPrice float64
			IName  string
			IData  string
		}
		err = query.Scan(&itemInfo).Error
		step6.TimeMs = time.Since(start).Seconds() * 1000
		if recordStep(tx, &step6, err, &response) {
			return
		}

		// 步骤7: 获取库存信息
		step7 := StepResult{StepName: fmt.Sprintf("get_stock_%d", olNumber)}
		start = time.Now()
		query = tx.Table("stock").
			Select("s_quantity, s_data, s_dist_01, s_dist_02, s_dist_03, s_dist_04, s_dist_05, "+
				"s_dist_06, s_dist_07, s_dist_08, s_dist_09, s_dist_10").
			Where("s_i_id = ? AND s_w_id = ?", item.OL_I_ID, item.OL_SUPPLY_W_ID)
		stmt = query.Session(&gorm.Session{DryRun: true}).Statement
		step7.SQL = replacePlaceholders(stmt.SQL.String(), stmt.Vars)

		// 2. 执行查询
		var stock struct {
			SQuantity int
			SData     string
			SDist01   string
			SDist02   string
			SDist03   string
			SDist04   string
			SDist05   string
			SDist06   string
			SDist07   string
			SDist08   string
			SDist09   string
			SDist10   string
		}
		err = query.Scan(&stock).Error

		step7.TimeMs = time.Since(start).Seconds() * 1000
		if recordStep(tx, &step7, err, &response) {
			return
		}

		// 步骤8: 更新库存
		step8 := StepResult{StepName: fmt.Sprintf("update_stock_%d", olNumber)}
		start = time.Now()
		newQuantity := stock.SQuantity - item.OL_QUANTITY
		if newQuantity < 0 {
			newQuantity += 91
		}

		sqlStr8 := `
			UPDATE stock SET s_quantity = ?
			WHERE s_i_id = ? AND s_w_id = ?`
		args8 := []interface{}{newQuantity, item.OL_I_ID, item.OL_SUPPLY_W_ID}
		step8.SQL = replacePlaceholders(sqlStr8, args8)

		// 2. 执行更新
		result := tx.Exec(sqlStr8, args8...)
		step8.TimeMs = time.Since(start).Seconds() * 1000
		if recordStep(tx, &step8, result.Error, &response) {
			return
		}

		// 步骤9: 计算订单行金额
		olAmount := float64(item.OL_QUANTITY) * itemInfo.IPrice *
			(1 + warehouse.WTax + district.DTax) * (1 - customer.CDiscount)

		// 步骤10: 插入订单行
		step10 := StepResult{StepName: fmt.Sprintf("insert_order_line_%d", olNumber)}
		start = time.Now()
		// 根据地区选择正确的dist_info
		distInfo := getDistInfo(stock, req.D_ID)
		if distInfo == "" {
			// 防止空字符串导致SQL语法错误
			distInfo = " "
		}

		sqlStr10 := `
			INSERT INTO order_line (ol_o_id, ol_d_id, ol_w_id, ol_number, 
				ol_i_id, ol_supply_w_id, ol_quantity, ol_amount, ol_dist_info)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		args10 := []interface{}{oID, req.D_ID, req.W_ID, olNumber,
			item.OL_I_ID, item.OL_SUPPLY_W_ID, item.OL_QUANTITY, olAmount, distInfo}
		step10.SQL = replacePlaceholders(sqlStr10, args10)

		// 2. 执行插入
		err = tx.Exec(sqlStr10, args10...).Error
		step10.TimeMs = time.Since(start).Seconds() * 1000
		if recordStep(tx, &step10, err, &response) {
			return
		}
		step10.TimeMs = time.Since(start).Seconds() * 1000
		if recordStep(tx, &step10, err, &response) {
			return
		}
	}

	// 提交事务
	stepCommit := StepResult{StepName: "commit"}
	start = time.Now()
	err = tx.Commit().Error
	stepCommit.TimeMs = time.Since(start).Seconds() * 1000
	recordStep(nil, &stepCommit, err, &response)

	if err == nil {
		response.Success = true
	}

	// 获取执行计划和分析索引使用
	analyzeSteps(&response)

	c.JSON(200, response)
}

// 辅助函数：记录步骤结果
func recordStep(tx *gorm.DB, step *StepResult, err error, response *NewOrderResponse) bool {
	if err != nil {
		step.TimeMs = 0
		response.Steps = append(response.Steps, *step)
		response.Success = false
		response.Message = err.Error()
		if tx != nil {
			tx.Rollback()
		}
		return true
	}
	response.Steps = append(response.Steps, *step)
	return false
}

// 辅助函数：获取地区信息
func getDistInfo(stock struct {
	SQuantity int
	SData     string
	SDist01   string
	SDist02   string
	SDist03   string
	SDist04   string
	SDist05   string
	SDist06   string
	SDist07   string
	SDist08   string
	SDist09   string
	SDist10   string
	// ... 其他字段
}, dID int) string {
	distFields := []string{
		stock.SDist01, stock.SDist02, stock.SDist03, stock.SDist04, stock.SDist05,
		stock.SDist06, stock.SDist07, stock.SDist08, stock.SDist09, stock.SDist10,
	}
	if dID >= 1 && dID <= 10 {
		return distFields[dID-1]
	}
	return distFields[0]
}

// 辅助函数：分析执行计划和索引使用
func analyzeSteps(response *NewOrderResponse) {
	for i := range response.Steps {
		step := &response.Steps[i]
		if step.SQL == "" {
			continue
		}

		var explainRows []map[string]interface{}
		result := global.DB.Raw("EXPLAIN FORMAT=JSON " + step.SQL).Scan(&explainRows)
		if result.Error != nil {
			step.ExplainPlan = result.Error.Error()
			continue
		}

		if len(explainRows) > 0 {
			explainJSON, _ := json.Marshal(explainRows[0])
			step.ExplainPlan = string(explainJSON)
			step.IndexUsed = extractIndexUsed(explainRows[0])
		}
	}
}

// 辅助函数：提取索引使用信息
func extractIndexUsed(explainRow map[string]interface{}) string {
	if key, ok := explainRow["key"]; ok && key != nil {
		return key.(string)
	}
	return "NONE"
}
func replacePlaceholders(sql string, args []interface{}) string {
	for _, arg := range args {
		strVal := ""
		switch v := arg.(type) {
		case string:
			// 空字符串处理为''
			if v == "" {
				strVal = "''"
			} else {
				strVal = "'" + v + "'"
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			strVal = fmt.Sprintf("%d", v)
		case float32, float64:
			strVal = fmt.Sprintf("%f", v)
		default:
			strVal = fmt.Sprintf("%v", arg)
		}
		sql = strings.Replace(sql, "?", strVal, 1)
	}
	return sql
}
