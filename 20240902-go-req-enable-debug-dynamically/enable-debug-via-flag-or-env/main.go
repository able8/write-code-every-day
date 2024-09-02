package main

import (
	"fmt"
	"log"
	"os"

	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
)

var globalClient = req.C()

var enableDebug bool

var uuidCmd = cobra.Command{
	Use: "uuid",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if enableDebug { // Enable debug mode if `--debug=true` or `DEBUG=true`.
			globalClient.EnableDumpAll()  // Dump all requests.
			globalClient.EnableDebugLog() // Output debug log.
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var result struct {
			Uuid string `json:"uuid"`
		}

		resp, err := globalClient.R().
			SetSuccessResult(&result). // Read uuid response into struct.
			Get("https://httpbin.org/uuid")

		if err != nil {
			return err
		}

		if resp.IsSuccessState() { // Print uuid returned by the API.
			fmt.Println(result.Uuid)
		} else {
			fmt.Println("bad response")
		}
		return nil
	},
}

func init() {
	uuidCmd.PersistentFlags().BoolVar(&enableDebug, "debug", os.Getenv("DEBUG") == "true", "Enable debug mode")
}

func main() {
	if err := uuidCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
