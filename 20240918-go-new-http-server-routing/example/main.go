package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

// https://github.com/vladimirvivien/go-httpmux-example/blob/main/main.go
// curl http://localhost:8080/tasks
// curl http://localhost:8080/tasks/one
// curl http://localhost:8080/tasks/create -d '{"id": "5", "description": "5 task", "completed": true}'

// Task item
type Task struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var (
	tasks = map[string]Task{
		"one":   {ID: "one", Description: "First task", Completed: true},
		"two":   {ID: "two", Description: "Second task", Completed: true},
		"three": {ID: "three", Description: "Third task", Completed: false},
		"four":  {ID: "four", Description: "Fourth task", Completed: true},
	}

	lock sync.RWMutex
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", getTasks)
	mux.HandleFunc("GET /tasks/{id}", getTask)
	mux.HandleFunc("POST /tasks/create", postTask)
	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", mux)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Error encountered", http.StatusInternalServerError)
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	lock.RLock()
	task, ok := tasks[idStr]
	lock.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Error encountered", http.StatusInternalServerError)
	}
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}
	id := fmt.Sprintf("%#X", rand.Intn(1024))
	task.ID = id

	lock.Lock()
	tasks[id] = task
	lock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Error encountered", http.StatusInternalServerError)
	}
}
