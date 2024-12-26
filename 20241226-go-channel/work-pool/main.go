package main

import (
	"fmt"
	"time"
)

func worker(id int, tasks <-chan int, results chan<- int) {
	for task := range tasks {
		fmt.Printf("Worker %d processing task %d\n", id, task)
		time.Sleep(time.Second) // 模拟任务处理时间
		results <- task * 2     // 返回处理结果
	}
}

func main() {
	tasks := make(chan int, 10)
	results := make(chan int, 10)

	for i := 1; i <= 3; i++ {
		go worker(i, tasks, results)
	}

	for j := 1; j <= 5; j++ {
		tasks <- j
	}
	close(tasks)

	for k := 1; k <= 5; k++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}
}
