package tpcc_api

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 请求参数结构体
type PaymentRequest struct {
	W_ID     int     `json:"w_id" binding:"required"`     // 仓库编号
	D_W_ID   int     `json:"d_w_id" binding:"required"`   // 地区仓库编号
	D_ID     int     `json:"d_id" binding:"required"`     // 地区编号
	C_W_ID   int     `json:"c_w_id" binding:"required"`   // 客户仓库编号
	C_D_ID   int     `json:"c_d_id" binding:"required"`   // 客户地区编号
	C_ID     int     `json:"c_id"`                        // 客户编号（可选）
	C_LAST   string  `json:"c_last"`                      // 客户姓氏（可选）
	H_AMOUNT float64 `json:"h_amount" binding:"required"` // 付款金额
}

// 步骤执行信息
type StepInfo struct {
	StepName  string        `json:"step_name"`  // 步骤名称
	TimeTaken time.Duration `json:"time_taken"` // 执行耗时
	Plan      string        `json:"plan"`       // 执行计划
	UsedIndex string        `json:"used_index"` // 使用的索引
}

// 事务响应结构
type PaymentResponse struct {
	Success bool       `json:"success"` // 是否成功
	Steps   []StepInfo `json:"steps"`   // 各步骤执行信息
	Message string     `json:"message"` // 附加消息
}

func (t *TpccApi) PaymentView(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBind(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	response := PaymentResponse{Steps: make([]StepInfo, 0)}
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			response.Success = false
			response.Message = fmt.Sprintf("Transaction panic: %v", r)
			res.OkWithData(response, c)
		}
	}()

	// 步骤1: 更新仓库余额
	step1Start := time.Now()
	if err := tx.Exec("UPDATE warehouse SET w_ytd = w_ytd + ? WHERE w_id = ?",
		req.H_AMOUNT, req.W_ID).Error; err != nil {
		tx.Rollback()
		res.FailWithError(err, c)
		return
	}
	step1Time := time.Since(step1Start)
	step1Plan, step1Index := getExplainInfo(tx, "UPDATE warehouse SET w_ytd = w_ytd + ? WHERE w_id = ?", req.H_AMOUNT, req.W_ID)
	response.Steps = append(response.Steps, StepInfo{
		StepName:  "update_warehouse",
		TimeTaken: step1Time,
		Plan:      step1Plan,
		UsedIndex: step1Index,
	})

	// 步骤2: 查询仓库信息
	type WarehouseInfo struct {
		WStreet1 string
		WStreet2 string
		WCity    string
		WState   string
		WZip     string
		WName    string
	}
	var warehouse WarehouseInfo
	step2Start := time.Now()
	if err := tx.Raw("SELECT w_street_1, w_street_2, w_city, w_state, w_zip, w_name FROM warehouse WHERE w_id = ?",
		req.W_ID).Scan(&warehouse).Error; err != nil {
		tx.Rollback()
		res.FailWithError(err, c)
		return
	}
	step2Time := time.Since(step2Start)
	step2Plan, step2Index := getExplainInfo(tx, "SELECT w_street_1, w_street_2, w_city, w_state, w_zip, w_name FROM warehouse WHERE w_id = ?", req.W_ID)
	response.Steps = append(response.Steps, StepInfo{
		StepName:  "select_warehouse",
		TimeTaken: step2Time,
		Plan:      step2Plan,
		UsedIndex: step2Index,
	})

	// 步骤3: 更新地区余额
	step3Start := time.Now()
	if err := tx.Exec("UPDATE district SET d_ytd = d_ytd + ? WHERE d_w_id = ? AND d_id = ?",
		req.H_AMOUNT, req.D_W_ID, req.D_ID).Error; err != nil {
		tx.Rollback()
		res.FailWithError(err, c)
		return
	}
	step3Time := time.Since(step3Start)
	step3Plan, step3Index := getExplainInfo(tx, "UPDATE district SET d_ytd = d_ytd + ? WHERE d_w_id = ? AND d_id = ?", req.H_AMOUNT, req.D_W_ID, req.D_ID)
	response.Steps = append(response.Steps, StepInfo{
		StepName:  "update_district",
		TimeTaken: step3Time,
		Plan:      step3Plan,
		UsedIndex: step3Index,
	})

	// 步骤4: 查询地区信息
	type DistrictInfo struct {
		DStreet1 string
		DStreet2 string
		DCity    string
		DState   string
		DZip     string
		DName    string
	}
	var district DistrictInfo
	step4Start := time.Now()
	if err := tx.Raw("SELECT d_street_1, d_street_2, d_city, d_state, d_zip, d_name FROM district WHERE d_w_id = ? AND d_id = ?",
		req.D_W_ID, req.D_ID).Scan(&district).Error; err != nil {
		tx.Rollback()
		res.FailWithError(err, c)
		return
	}
	step4Time := time.Since(step4Start)
	step4Plan, step4Index := getExplainInfo(tx, "SELECT d_street_1, d_street_2, d_city, d_state, d_zip, d_name FROM district WHERE d_w_id = ? AND d_id = ?", req.D_W_ID, req.D_ID)
	response.Steps = append(response.Steps, StepInfo{
		StepName:  "select_district",
		TimeTaken: step4Time,
		Plan:      step4Plan,
		UsedIndex: step4Index,
	})

	// 步骤5: 根据客户标识类型处理客户信息
	type CustomerInfo struct {
		CId        uint
		CFirst     string
		CMiddle    string
		CLast      string
		CStreet1   string
		CStreet2   string
		CCity      string
		CState     string
		CZip       string
		CPhone     string
		CCredit    string
		CCreditLim float64
		CDiscount  float64
		CBalance   float64
		CSince     time.Time
		CData      string
	}
	var customer CustomerInfo

	if req.C_LAST != "" {
		// 5.1 按姓氏查询客户
		var namecnt int
		step5_1Start := time.Now()
		if err := tx.Raw("SELECT COUNT(c_id) FROM customer WHERE c_last = ? AND c_d_id = ? AND c_w_id = ?",
			req.C_LAST, req.C_D_ID, req.C_W_ID).Scan(&namecnt).Error; err != nil {
			tx.Rollback()
			res.FailWithError(err, c)
			return
		}
		step5_1Time := time.Since(step5_1Start)
		step5_1Plan, step5_1Index := getExplainInfo(tx, "SELECT COUNT(c_id) FROM customer WHERE c_last = ? AND c_d_id = ? AND c_w_id = ?", req.C_LAST, req.C_D_ID, req.C_W_ID)
		response.Steps = append(response.Steps, StepInfo{
			StepName:  "count_customers_by_last_name",
			TimeTaken: step5_1Time,
			Plan:      step5_1Plan,
			UsedIndex: step5_1Index,
		})

		// 计算中间客户位置
		offset := (namecnt + 1) / 2
		if offset > 0 {
			offset--
		}

		step5_2Start := time.Now()
		query := `SELECT c_first, c_middle, c_id, c_street_1, c_street_2, c_city, c_state, c_zip, 
                 c_phone, c_credit, c_credit_lim, c_discount, c_balance, c_since
                 FROM customer 
                 WHERE c_w_id = ? AND c_d_id = ? AND c_last = ? 
                 ORDER BY c_first 
                 LIMIT 1 OFFSET ?`
		if err := tx.Raw(query, req.C_W_ID, req.C_D_ID, req.C_LAST, offset).Scan(&customer).Error; err != nil {
			tx.Rollback()
			res.FailWithError(err, c)
			return
		}
		step5_2Time := time.Since(step5_2Start)
		step5_2Plan, step5_2Index := getExplainInfo(tx, query, req.C_W_ID, req.C_D_ID, req.C_LAST, offset)
		response.Steps = append(response.Steps, StepInfo{
			StepName:  "select_customer_by_last_name",
			TimeTaken: step5_2Time,
			Plan:      step5_2Plan,
			UsedIndex: step5_2Index,
		})
	} else {
		// 5.2 按ID查询客户
		step5Start := time.Now()
		query := `SELECT c_first, c_middle, c_last, c_street_1, c_street_2, c_city, c_state, c_zip, 
                 c_phone, c_credit, c_credit_lim, c_discount, c_balance, c_since
                 FROM customer 
                 WHERE c_w_id = ? AND c_d_id = ? AND c_id = ?`
		if err := tx.Raw(query, req.C_W_ID, req.C_D_ID, req.C_ID).Scan(&customer).Error; err != nil {
			tx.Rollback()
			res.FailWithError(err, c)
			return
		}
		step5Time := time.Since(step5Start)
		step5Plan, step5Index := getExplainInfo(tx, query, req.C_W_ID, req.C_D_ID, req.C_ID)
		response.Steps = append(response.Steps, StepInfo{
			StepName:  "select_customer_by_id",
			TimeTaken: step5Time,
			Plan:      step5Plan,
			UsedIndex: step5Index,
		})
	}

	// 更新客户余额
	customer.CBalance += req.H_AMOUNT

	// 步骤6: 根据信用类型更新客户
	if customer.CCredit == "BC" {
		step6_1Start := time.Now()
		if err := tx.Raw("SELECT c_data FROM customer WHERE c_w_id = ? AND c_d_id = ? AND c_id = ?",
			req.C_W_ID, req.C_D_ID, customer.CId).Scan(&customer.CData).Error; err != nil {
			tx.Rollback()
			res.FailWithError(err, c)
			return
		}
		step6_1Time := time.Since(step6_1Start)
		step6_1Plan, step6_1Index := getExplainInfo(tx, "SELECT c_data FROM customer WHERE c_w_id = ? AND c_d_id = ? AND c_id = ?", req.C_W_ID, req.C_D_ID, customer.CId)
		response.Steps = append(response.Steps, StepInfo{
			StepName:  "select_customer_data",
			TimeTaken: step6_1Time,
			Plan:      step6_1Plan,
			UsedIndex: step6_1Index,
		})

		// 生成新客户数据
		newData := fmt.Sprintf("| %4d %2d %4d %2d %4d $%7.2f %s %s",
			customer.CId, req.C_D_ID, req.C_W_ID, req.D_ID, req.W_ID, req.H_AMOUNT,
			time.Now().Format("2006-01-02 15:04:05"), customer.CData)
		if len(newData) > 500 {
			newData = newData[:500]
		}

		step6_2Start := time.Now()
		if err := tx.Exec("UPDATE customer SET c_balance = ?, c_data = ? WHERE c_w_id = ? AND c_d_id = ? AND c_id = ?",
			customer.CBalance, newData, req.C_W_ID, req.C_D_ID, customer.CId).Error; err != nil {
			tx.Rollback()
			res.FailWithError(err, c)
			return
		}
		step6_2Time := time.Since(step6_2Start)
		step6_2Plan, step6_2Index := getExplainInfo(tx, "UPDATE customer SET c_balance = ?, c_data = ? WHERE c_w_id = ? AND c_d_id = ? AND c_id = ?", customer.CBalance, newData, req.C_W_ID, req.C_D_ID, customer.CId)
		response.Steps = append(response.Steps, StepInfo{
			StepName:  "update_customer_with_data",
			TimeTaken: step6_2Time,
			Plan:      step6_2Plan,
			UsedIndex: step6_2Index,
		})
	} else {
		step6Start := time.Now()
		if err := tx.Exec("UPDATE customer SET c_balance = ? WHERE c_w_id = ? AND c_d_id = ? AND c_id = ?",
			customer.CBalance, req.C_W_ID, req.C_D_ID, customer.CId).Error; err != nil {
			tx.Rollback()
			res.FailWithError(err, c)
			return
		}
		step6Time := time.Since(step6Start)
		step6Plan, step6Index := getExplainInfo(tx, "UPDATE customer SET c_balance = ? WHERE c_w_id = ? AND c_d_id = ? AND c_id = ?", customer.CBalance, req.C_W_ID, req.C_D_ID, customer.CId)
		response.Steps = append(response.Steps, StepInfo{
			StepName:  "update_customer",
			TimeTaken: step6Time,
			Plan:      step6Plan,
			UsedIndex: step6Index,
		})
	}

	// 步骤7: 插入历史记录
	step7Start := time.Now()
	hData := fmt.Sprintf("%-10s%-10s    ", warehouse.WName, district.DName) // 组合仓库和地区名称
	if err := tx.Exec(`INSERT INTO history 
        (h_c_d_id, h_c_w_id, h_c_id, h_d_id, h_w_id, h_date, h_amount, h_data) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.C_D_ID, req.C_W_ID, customer.CId, req.D_ID, req.W_ID, time.Now(), req.H_AMOUNT, hData).Error; err != nil {
		tx.Rollback()
		res.FailWithError(err, c)
		return
	}
	step7Time := time.Since(step7Start)
	step7Plan, step7Index := getExplainInfo(tx, `INSERT INTO history 
        (h_c_d_id, h_c_w_id, h_c_id, h_d_id, h_w_id, h_date, h_amount, h_data) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.C_D_ID, req.C_W_ID, customer.CId, req.D_ID, req.W_ID, time.Now(), req.H_AMOUNT, hData)
	response.Steps = append(response.Steps, StepInfo{
		StepName:  "insert_history",
		TimeTaken: step7Time,
		Plan:      step7Plan,
		UsedIndex: step7Index,
	})

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		response.Success = false
		response.Message = "Commit failed: " + err.Error()
		res.OkWithData(response, c)
		return
	}

	response.Success = true
	response.Message = "Payment transaction completed successfully"
	res.OkWithData(response, c)
}

// 获取SQL执行计划和索引使用情况
func getExplainInfo(db *gorm.DB, query string, args ...interface{}) (string, string) {
	var explainRows []map[string]interface{}
	explainQuery := "EXPLAIN " + query
	if err := db.Raw(explainQuery, args...).Scan(&explainRows).Error; err != nil {
		return "Explain error", ""
	}

	var planBuilder strings.Builder
	var usedIndex string

	for _, row := range explainRows {
		if idx, ok := row["key"]; ok && idx != nil {
			usedIndex = idx.(string)
		}
		for k, v := range row {
			planBuilder.WriteString(fmt.Sprintf("%s: %v | ", k, v))
		}
		planBuilder.WriteString("\n")
	}

	return planBuilder.String(), usedIndex
}
