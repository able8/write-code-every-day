package main

import (
	"fmt"
	"time"
)

func broadcast(channel <-chan string, id int) {
	for msg := range channel {
		fmt.Printf("Listener %d received message: %s\n", id, msg)
	}
}

func main() {
	eventChannel := make(chan string)

	for i := 1; i <= 3; i++ {
		go broadcast(eventChannel, i)
	}

	messages := []string{"Event 1", "Event 2", "Event 3"}
	for _, msg := range messages {
		eventChannel <- msg
		time.Sleep(500 * time.Millisecond)
	}
	close(eventChannel)
}
