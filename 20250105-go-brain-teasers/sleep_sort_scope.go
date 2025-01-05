package main

import (
	"fmt"
	"sync"
	"time"
)
func main() {
	var wg sync.WaitGroup
	for _, n := range []int{3, 1, 2} {
		n := n // <1>
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(n) * time.Millisecond)
			fmt.Printf("%d ", n)
		}()
	}
	wg.Wait()
	fmt.Println()
}
