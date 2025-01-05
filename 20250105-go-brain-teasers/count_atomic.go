package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var count int64 // Atomic works with int64, not int
	var wg sync.WaitGroup

	for i := 0; i < 1_000_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&count, 1)
		}()

	}
	wg.Wait()
	fmt.Println(count)
}
