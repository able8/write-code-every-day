package main

import (
	"fmt"
	"sync"
	"time"
)

// Helper function to run a function in a goroutine and manage the wait group
func runInGoroutine(wg *sync.WaitGroup, f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
}

func main() {
	var wg sync.WaitGroup

	// Define a few functions to run in goroutines
	task1 := func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Task 1 - Count:", i)
			time.Sleep(100 * time.Millisecond) // Simulate work
		}
	}

	task2 := func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Task 2 - Count:", i)
			time.Sleep(150 * time.Millisecond) // Simulate work
		}
	}

	task3 := func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Task 3 - Count:", i)
			time.Sleep(200 * time.Millisecond) // Simulate work
		}
		panic("test")
	}

	// Run tasks in goroutines
	runInGoroutine(&wg, task1)
	runInGoroutine(&wg, task2)
	runInGoroutine(&wg, task3)

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All tasks completed.")
}
