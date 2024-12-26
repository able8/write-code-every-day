package main

import (
	"fmt"
	"time"
)

type LogMessage struct {
	Level   string // 日志级别
	Message string // 日志内容
}

// 日志处理器
func logProcessor(logChannel <-chan LogMessage, batchSize int, flushInterval time.Duration) {
	var logBatch []LogMessage
	timer := time.NewTimer(flushInterval)

	for {
		select {
		case log := <-logChannel:
			// 将日志添加到当前批次
			logBatch = append(logBatch, log)

			// 当达到批量大小时处理日志
			if len(logBatch) >= batchSize {
				processBatch(logBatch)     // 处理日志批次
				logBatch = nil             // 清空批次
				timer.Reset(flushInterval) // 重置定时器
			}

		case <-timer.C:
			// 定时器触发时处理未满批次的日志
			if len(logBatch) > 0 {
				processBatch(logBatch)
				logBatch = nil
			}
			timer.Reset(flushInterval)
		}
	}
}

// 处理日志批次的函数（伪代码）
func processBatch(logs []LogMessage) {
	for _, log := range logs {
		fmt.Printf("[%s] %s\n", log.Level, log.Message)
	}
}

func main() {
	logChannel := make(chan LogMessage, 100) // 创建带缓冲的 channel
	batchSize := 10                          // 设置批量大小
	flushInterval := 5 * time.Second         // 设置定时刷新间隔

	go logProcessor(logChannel, batchSize, flushInterval)

	// 模拟日志生成
	for i := 0; i < 50; i++ {
		logChannel <- LogMessage{
			Level:   "INFO",
			Message: fmt.Sprintf("Log message %d", i),
		}
		time.Sleep(200 * time.Millisecond) // 模拟日志生成间隔
	}

	// 确保主 goroutine 不会过早退出
	time.Sleep(10 * time.Second)
}
