package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 模拟清理资源的函数
func cleanUp() {
	fmt.Println("Cleaning up resources...")
	// 模拟清理任务，如关闭数据库连接、清理缓存、保存日志等
	time.Sleep(2 * time.Second) // 假设清理任务需要 2 秒钟
	fmt.Println("Resources cleaned up.")
}

func main() {
	// 创建一个取消的上下文，用于控制优雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建一个信号通道，用于接收操作系统的信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM) // 捕获 SIGINT 和 SIGTERM 信号

	// 启动一个 goroutine 进行信号监听
	go func() {
		sig := <-signalChan
		fmt.Println("Received signal:", sig)
		// 收到信号后取消上下文，进行清理
		cancel()
	}()

	// 模拟主程序运行
	fmt.Println("Application started.")
	for {
		select {
		case <-ctx.Done():
			// 收到关闭信号，执行清理
			cleanUp()
			fmt.Println("Shutting down application...")
			return
		default:
			// 模拟应用程序工作
			time.Sleep(1 * time.Second)
			fmt.Println("working...")
		}
	}
}
