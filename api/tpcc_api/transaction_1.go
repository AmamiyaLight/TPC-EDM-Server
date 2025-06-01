package tpcc_api

import (
	"TPC-EDM-Server/global"
	"database/sql"
	"encoding/json"
	"fmt"
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
	StepName     string  `json:"step_name"`
	SQL          string  `json:"sql"`
	TimeMs       float64 `json:"time_ms"`
	ExplainPlan  string  `json:"explain_plan"`
	IndexUsed    string  `json:"index_used"`
	RowsAffected int64   `json:"rows_affected"`
}

type NewOrderResponse struct {
	Success bool         `json:"success"`
	Total   float64      `json:"total"`
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
	var total float64
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
	var result CustomerWarehouseResult

	// 使用显式别名修复冲突
	err := tx.Table("customer AS c").
		Select("c.c_discount, c.c_last, c.c_credit, w.w_tax").
		Joins("JOIN warehouse AS w ON w.w_id = c.c_w_id").
		Where("c.c_w_id = ? AND c.c_d_id = ? AND c.c_id = ?", req.C_W_ID, req.C_D_ID, req.C_ID).
		Where("w.w_id = ?", req.W_ID).
		Scan(&result).Error

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
	var district struct {
		DNextOID int
		DTax     float64
	}
	err = tx.Table("district").
		Select("d_next_o_id, d_tax").
		Where("d_id = ? AND d_w_id = ?", req.D_ID, req.D_W_ID).
		Scan(&district).Error
	step2.TimeMs = time.Since(start).Seconds() * 1000
	if recordStep(tx, &step2, err, &response) {
		return
	}

	// 步骤3: 更新地区订单ID
	step3 := StepResult{StepName: "update_district"}
	start = time.Now()
	err = tx.Exec("UPDATE district SET d_next_o_id = ? WHERE d_id = ? AND d_w_id = ?",
		district.DNextOID+1, req.D_ID, req.D_W_ID).Error
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

	err = tx.Exec(`
		INSERT INTO orders (o_id, o_d_id, o_w_id, o_c_id, o_entry_d, o_ol_cnt, o_all_local)
		VALUES (?, ?, ?, ?, NOW(), ?, ?)`,
		oID, req.D_ID, req.W_ID, req.C_ID, req.O_OL_CNT, oAllLocal).Error
	step4.TimeMs = time.Since(start).Seconds() * 1000
	if recordStep(tx, &step4, err, &response) {
		return
	}

	// 步骤5: 插入新订单
	step5 := StepResult{StepName: "insert_new_order"}
	start = time.Now()
	err = tx.Exec(`
		INSERT INTO new_orders (no_o_id, no_d_id, no_w_id)
		VALUES (?, ?, ?)`,
		oID, req.D_ID, req.W_ID).Error
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
		var itemInfo struct {
			IPrice float64
			IName  string
			IData  string
		}
		err = tx.Table("item").
			Select("i_price, i_name, i_data").
			Where("i_id = ?", item.OL_I_ID).
			Scan(&itemInfo).Error
		step6.TimeMs = time.Since(start).Seconds() * 1000
		if recordStep(tx, &step6, err, &response) {
			return
		}

		// 步骤7: 获取库存信息
		step7 := StepResult{StepName: fmt.Sprintf("get_stock_%d", olNumber)}
		start = time.Now()
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
		err = tx.Table("stock").
			Select("s_quantity, s_data, s_dist_01, s_dist_02, s_dist_03, s_dist_04, s_dist_05, "+
				"s_dist_06, s_dist_07, s_dist_08, s_dist_09, s_dist_10").
			Where("s_i_id = ? AND s_w_id = ?", item.OL_I_ID, item.OL_SUPPLY_W_ID).
			Scan(&stock).Error
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

		result := tx.Exec(`
			UPDATE stock SET s_quantity = ?
			WHERE s_i_id = ? AND s_w_id = ?`,
			newQuantity, item.OL_I_ID, item.OL_SUPPLY_W_ID)
		step8.RowsAffected = result.RowsAffected
		step8.TimeMs = time.Since(start).Seconds() * 1000
		if recordStep(tx, &step8, result.Error, &response) {
			return
		}

		// 步骤9: 计算订单行金额
		olAmount := float64(item.OL_QUANTITY) * itemInfo.IPrice *
			(1 + warehouse.WTax + district.DTax) * (1 - customer.CDiscount)
		total += olAmount

		// 步骤10: 插入订单行
		step10 := StepResult{StepName: fmt.Sprintf("insert_order_line_%d", olNumber)}
		start = time.Now()
		// 根据地区选择正确的dist_info
		distInfo := getDistInfo(stock, req.D_ID)

		err = tx.Exec(`
			INSERT INTO order_line (ol_o_id, ol_d_id, ol_w_id, ol_number, 
				ol_i_id, ol_supply_w_id, ol_quantity, ol_amount, ol_dist_info)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			oID, req.D_ID, req.W_ID, olNumber,
			item.OL_I_ID, item.OL_SUPPLY_W_ID, item.OL_QUANTITY, olAmount, distInfo).Error
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
		response.Total = total
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
