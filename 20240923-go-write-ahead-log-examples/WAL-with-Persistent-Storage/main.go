package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type WAL struct {
	mu         sync.Mutex
	log        []string
	commit     bool
	persistent string
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

	// Save the log to persistent storage
	if err := saveLog(w.persistent, w.log); err != nil {
		fmt.Println("Error saving WAL:", err)
	}
}

func saveLog(file string, log []string) error {
	data, err := json.Marshal(log)
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}

func main() {
	wal := &WAL{}
	wal.persistent = "wal.json"

	// Append some messages to the log
	wal.Append("Message 1")
	wal.Append("Message 2")

	// Commit the log
	wal.Commit()
}
