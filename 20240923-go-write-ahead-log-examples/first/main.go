package main

import (
	"fmt"
	"sync"
)

type WAL struct {
	mu     sync.Mutex
	log    []string
	commit bool
}

func (w *WAL) Append(msg string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.log = append(w.log, msg)
}

func (w *WAL) Commit() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.commit = true
	fmt.Println("Committing WAL:", w.log)
}

func main() {
	wal := &WAL{}

	// Append some messages to the log
	wal.Append("Message 1")
	wal.Append("Message 2")

	// Commit the log
	wal.Commit()
}
