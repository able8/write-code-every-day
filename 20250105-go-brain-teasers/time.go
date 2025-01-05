package main

import (
	"fmt"
	"time"
)

func main() {
	timeout := 3
	fmt.Printf("before ")
	time.Sleep(timeout * time.Millisecond)
	time.Sleep(3 * time.Millisecond)
	time.Sleep(time.Duration(timeout) * time.Second)
	fmt.Println("after")
}
