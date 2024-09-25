package main

import (
	"fmt"
	"sync"
	"time"

	"math/rand"
)

func main() {
	a := asChan(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	b := asChan(10, 11, 12, 13, 14, 15, 16, 17, 18, 19)
	c := asChan(20, 21, 22, 23, 24, 25, 26, 27, 28, 29)
	for v := range merge(a, b, c) {
		fmt.Println(v)
	}
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(cs))

	for _, c := range cs {

		// Fixing For Loops in Go 1.22
		go func() {
			for v := range c {
				out <- v
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)

	}()

	return out
}

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}
