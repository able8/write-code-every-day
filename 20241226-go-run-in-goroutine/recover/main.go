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

		// Recover from panic if f() panics
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in runInGoroutine:", r)
			}
		}()

		f()
	}()
}

func main() {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 2) // Example: limit to 2 concurrent goroutines

	for i := 0; i < 10; i++ {
		runInGoroutine(&wg, sem, func() {
			time.Sleep(1 * time.Second)
			fmt.Println("Running function")
			panic("test")
		})
	}

	wg.Wait() // Wait for all goroutines to finish
}
