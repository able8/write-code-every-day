package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func worker(id int, ctx context.Context) {
	fmt.Printf("Worker %d started\n", id)
	for {
		select {
		case <-ctx.Done():
			// 收到取消信号，优雅退出
			fmt.Printf("Worker %d is stopping\n", id)
			return
		default:
			// 模拟执行工作任务
			time.Sleep(1 * time.Second)
			fmt.Printf("Worker %d is working...\n", id)
		}
	}
}

func main() {
	// 创建一个带取消的上下文，用于优雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建信号通道，用于捕获系统信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动多个工作 goroutine
	for i := 1; i <= 3; i++ {
		go worker(i, ctx)
	}

	// 等待终止信号
	sig := <-signalChan
	fmt.Println("Received signal:", sig)

	// 收到信号后，取消上下文，所有 goroutine 会响应并退出
	cancel()

	// 等待所有 goroutine 完成
	time.Sleep(3 * time.Second) // 给予足够的时间完成清理工作
	fmt.Println("Application shut down gracefully.")
}
