package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", ExampleHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

	// Double check it's a post request being made
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	// Must call ParseForm() before working with data
	r.ParseForm()

	// Log all data. Form is a map[]
	log.Println(r.Form)

	// Print the data back. We can use Form.Get() or Form["name"][0]
	fmt.Fprintf(w, "Hello "+r.Form.Get("name"))
}

// curl -X POST localhost:8080 -d "name=able"
