package main

import (
	"fmt"
	"net/http"
)

func main() {

	// ServeMux is an HTTP request multiplexer.
	// It matches the URL of each incoming request against a list of registered
	// patterns and calls the handler for the pattern that
	// most closely matches the URL.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /path/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "got path, ", r.URL)
	})

	mux.HandleFunc("/task/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "handling task with id=%v\n", id)
	})

	mux.HandleFunc("/task/{id}/status/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "handling task status with id=%v\n", id)
	})
	mux.HandleFunc("/task/0/{action}/", func(w http.ResponseWriter, r *http.Request) {
		action := r.PathValue("action")
		fmt.Fprintf(w, "handling task 0 with action=%v\n", action)
	})

	http.ListenAndServe("localhost:8090", mux)
}
