package file

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"time"
)

func ProcessFileInsert[T any](
	c *gin.Context,
	parseFunc func(string) (T, error), // 行数据解析函数
	batchSize int, // 批次大小
) (int, error) {
	// 文件处理部分
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return 0, errors.New("文件获取失败")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return 0, errors.New("文件打开失败")
	}
	defer file.Close()

	// 扫描器初始化
	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 1024*1024)

	// 初始化批次
	batch := make([]T, 0, batchSize)
	total := 0

	// 处理每一行
	for scanner.Scan() {
		item, err := parseFunc(scanner.Text())
		if err != nil {
			continue // 或者根据需求处理错误
		}

		batch = append(batch, item)
		total++

		// 批量插入
		if len(batch) >= batchSize {
			if err := global.DB.Transaction(func(tx *gorm.DB) error {
				return tx.CreateInBatches(batch, len(batch)).Error
			}); err != nil {
				return total, err
			}
			batch = batch[:0]
		}
	}

	// 处理剩余数据
	if len(batch) > 0 {
		if err := global.DB.Transaction(func(tx *gorm.DB) error {
			return tx.CreateInBatches(batch, len(batch)).Error
		}); err != nil {
			return total, err
		}
	}

	return total, scanner.Err()
}

type DownloadProgress struct {
	Total     int64
	Processed int64
	StartTime time.Time
}

func DownloadHandler[T any](
	c *gin.Context,
	filename string,
	headers []string,
	getTotalFunc func() int64,
	fetchPageFunc func(offset, limit int) ([]T, error),
	convertFunc func(item T) []string,
) {
	startTime := time.Now()
	total := getTotalFunc()
	progress := &DownloadProgress{
		Total:     total,
		Processed: 0,
		StartTime: startTime,
	}

	// 设置响应头
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Transfer-Encoding", "chunked")

	// 创建CSV写入器
	writer := csv.NewWriter(c.Writer)
	writer.Comma = ','
	defer writer.Flush()

	// 写入CSV头部
	if err := writer.Write(headers); err != nil {
		res.FailWithError(err, c)
		return
	}
	c.Writer.Flush()

	// 分页获取数据并写入
	pageSize := 1000
	for offset := 0; offset < int(total); {
		items, err := fetchPageFunc(offset, pageSize)
		if err != nil {
			res.FailWithError(err, c)
			return
		}

		// 转换并写入数据
		for _, item := range items {
			record := convertFunc(item)
			if err := writer.Write(record); err != nil {
				res.FailWithError(err, c)
				return
			}
			progress.Processed++
		}

		// 刷新缓冲并更新进度
		writer.Flush()
		c.Writer.Flush()
		offset += len(items)

		// 实时进度日志（实际应用中可推送到前端）
		logProgress(progress)
	}

	// 完成日志
	duration := time.Since(startTime)
	log.Printf("下载完成: %s, 记录数: %d, 耗时: %v", filename, total, duration)
}
func logProgress(p *DownloadProgress) {
	elapsed := time.Since(p.StartTime)
	percent := float64(p.Processed) / float64(p.Total) * 100
	rate := float64(p.Processed) / elapsed.Seconds()

	log.Printf("下载进度: %.1f%% | 已处理: %d/%d | 速率: %.1f rec/s | 耗时: %v",
		percent, p.Processed, p.Total, rate, elapsed.Round(time.Second))
}
