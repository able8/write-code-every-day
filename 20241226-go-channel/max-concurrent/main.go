package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int) {
	fmt.Printf("Worker %d is working...\n", id)
	time.Sleep(1 * time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	const maxConcurrent = 3
	semaphore := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() { <-semaphore }()
			worker(id)
		}(i)
	}

	wg.Wait()
}
