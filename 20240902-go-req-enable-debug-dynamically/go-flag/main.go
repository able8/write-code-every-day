package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/imroc/req/v3"
)

var globalClient = req.C()

var enableDebug bool

func init() {
	flag.BoolVar(&enableDebug, "debug", os.Getenv("DEBUG") == "true", "Enable debug mode")
}

func main() {
	flag.Parse()
	if enableDebug { // Enable debug mode if `--enableDebug=true` or `DEBUG=true`.
		globalClient.EnableDumpAll()  // Dump all requests.
		globalClient.EnableDebugLog() // Output debug log.
	}
	var result struct {
		Uuid string `json:"uuid"`
	}

	resp, err := globalClient.R().
		SetSuccessResult(&result). // Read uuid response into struct.
		Get("https://httpbin.org/uuid")

	if err != nil {
		panic(err)
	}

	if resp.IsSuccessState() { // Print uuid returned by the API.
		fmt.Println(result.Uuid)
	} else {
		fmt.Println("bad response")
	}
}
