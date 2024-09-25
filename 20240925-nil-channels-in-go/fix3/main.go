package main

import (
	"fmt"
	"log"
	"time"

	"math/rand"
)

// https://medium.com/justforfunc/why-are-there-nil-channels-in-go-9877cc0b2308

// write a function that given two channels a and b of some type
// returns one channel c of the same type.
// Every element received in a or b will be sent to c,
// and once both a and b are closed c will be closed too.

func main() {
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Println(v)
	}
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {

		defer close(c)

		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					log.Println("a is done")
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					log.Println("b is done")
					// bdone = true
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
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
