package file

import (
	"TPC-EDM-Server/global"
	"bufio"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
