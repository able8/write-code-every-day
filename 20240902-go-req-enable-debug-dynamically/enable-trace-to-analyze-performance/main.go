package main

import (
	"fmt"
	"log"

	"github.com/imroc/req/v3"
)

func main() {
	// Enable trace at request level
	client := req.C()
	resp, err := client.R().EnableTrace().Get("https://api.github.com/users/imroc")
	if err != nil {
		log.Fatal(err)
	}
	trace := resp.TraceInfo()  // Use `resp.Request.TraceInfo()` to avoid unnecessary struct copy in production.
	fmt.Println(trace.Blame()) // Print out exactly where the http request is slowing down.
	fmt.Println("----------")
	fmt.Println(trace) // Print details

	resp, err = client.R().EnableTrace().Get("https://api.github.com/users/imroc")
	if err != nil {
		log.Fatal(err)
	}
	trace = resp.TraceInfo()   // Use `resp.Request.TraceInfo()` to avoid unnecessary struct copy in production.
	fmt.Println(trace.Blame()) // Print out exactly where the http request is slowing down.
	fmt.Println("----------")
	fmt.Println(trace) // Print details

}
