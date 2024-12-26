package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			ch <- i // 发送数据到 channel
		}
		close(ch) // 关闭 channel
	}()

	for val := range ch { // 从 channel 中读取数据
		fmt.Println(val)
	}
}
