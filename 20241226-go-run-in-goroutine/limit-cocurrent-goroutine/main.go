package main

import (
	"fmt"
	"sync"
	"time"
)

// Helper function to run a function in a goroutine and manage the wait group
func runInGoroutine(wg *sync.WaitGroup, sem chan struct{}, f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		sem <- struct{}{}        // Acquire a token
		defer func() { <-sem }() // Release the token when done
		f()
	}()
}

func main() {
	var wg sync.WaitGroup
	const maxConcurrentGoroutines = 2                   // Limit the number of concurrent goroutines
	sem := make(chan struct{}, maxConcurrentGoroutines) // Create a buffered channel

	// Define a few functions to run in goroutines
	task1 := func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Task 1 - Count:", i)
			time.Sleep(500 * time.Millisecond) // Simulate work
		}
	}

	task2 := func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Task 2 - Count:", i)
			time.Sleep(550 * time.Millisecond) // Simulate work
		}
	}

	task3 := func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Task 3 - Count:", i)
			time.Sleep(500 * time.Millisecond) // Simulate work
		}
	}

	// Run tasks in goroutines
	runInGoroutine(&wg, sem, task1)
	runInGoroutine(&wg, sem, task2)
	runInGoroutine(&wg, sem, task3)

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All tasks completed.")
}
