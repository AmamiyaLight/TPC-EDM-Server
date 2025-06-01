package service

import (
	"TPC-EDM-Server/global"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
	"time"
)

func TpchTest() {
	// 1. 读取SQL文件
	sqlContent, err := os.ReadFile("./service/tpch-stream.sql")
	if err != nil {
		logrus.Errorf("读取SQL文件失败: %v", err)
		return
	}

	// 2. 分割SQL语句
	queries := parseSQLQueries(string(sqlContent))
	if len(queries) == 0 {
		logrus.Warning("未找到有效的SQL语句")
		return
	}

	// 3. 准备性能统计
	var (
		wg            sync.WaitGroup
		mu            sync.Mutex
		totalQueries  int
		totalLatency  time.Duration
		workerResults = make([]workerStats, 4)
	)

	// 4. 启动4个worker
	startTime := time.Now()
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			stats := executeQueries(workerID, queries)
			mu.Lock()
			defer mu.Unlock()
			workerResults[workerID] = stats
			totalQueries += stats.queryCount
			totalLatency += stats.totalLatency
		}(i)
	}
	wg.Wait()
	totalDuration := time.Since(startTime)

	// 5. 计算性能指标
	throughput := float64(totalQueries) / totalDuration.Seconds()
	avgLatency := totalLatency.Seconds() / float64(totalQueries) * 1000 // 毫秒

	// 6. 打印结果
	logrus.Infof("============== 性能报告 ==============")
	logrus.Infof("总执行时间: %.4f 秒", totalDuration.Seconds())
	logrus.Infof("总查询数量: %d", totalQueries)
	logrus.Infof("系统吞吐量: %.2f 查询/秒", throughput)
	logrus.Infof("平均延迟: %.2f 毫秒", avgLatency)
	logrus.Infof("------------------------------------")

	for i, stats := range workerResults {
		logrus.Infof("Worker %d: 查询数=%d 总耗时=%.4fs 平均延迟=%.2fms",
			i, stats.queryCount, stats.totalLatency.Seconds(),
			stats.totalLatency.Seconds()/float64(stats.queryCount)*1000)
	}
	logrus.Infof("====================================")
}

// 解析SQL文件内容
func parseSQLQueries(content string) []string {
	rawQueries := strings.Split(content, ";")
	queries := make([]string, 0, len(rawQueries))

	for _, q := range rawQueries {
		cleaned := strings.TrimSpace(q)
		if cleaned != "" {
			queries = append(queries, cleaned)
		}
	}
	return queries
}

// Worker统计结构
type workerStats struct {
	queryCount   int
	totalLatency time.Duration
}

// 执行查询序列
func executeQueries(workerID int, queries []string) workerStats {
	var stats workerStats
	for _, query := range queries {
		start := time.Now()
		result := global.DB.Exec(query)
		latency := time.Since(start)

		stats.queryCount++
		stats.totalLatency += latency

		if result.Error != nil {
			logrus.Warnf("[Worker %d] 查询执行失败: %v\nSQL: %s", workerID, result.Error, query)
		} else {
			logrus.Debugf("[Worker %d] 执行成功 (%.2fms): %s",
				workerID, latency.Seconds()*1000, query)
		}
	}
	return stats
}
