package main

import (
	"encoding/binary"
	"flag"
	"fmt"

	"github.com/imroc/req/v3"
)

// go run main.go -d '{"id":1,"role":"ADMIN"}' localhost:10000/example.UserService/AddUser
// go run main.go -d '{}' localhost:10000/example.UserService/ListUsers

func main() {
	dataPtr := flag.String("d", "", "JSON payload data")

	flag.Parse()

	args := flag.Args()
	if *dataPtr == "" || len(args) == 0 {
		fmt.Println("Error: The -d flag and a URL argument are required.")
		flag.Usage()
		return
	}

	url := args[len(args)-1] // The last argument is the URL

	data, err := createRequestData([]byte(*dataPtr))
	if err != nil {
		fmt.Println("Error:", err)
	}

	client := req.C().DevMode().EnableForceHTTP2().EnableInsecureSkipVerify()

	response, err := client.R().
		SetHeader("Content-Type", "application/grpc+json").
		SetHeader("TE", "trailers").
		SetBody(data).
		Post("https://" + url)

	if err != nil {
		fmt.Println("Request error:", err)
		return
	}

	fmt.Println("Response:", response)
}

// createRequestData prepares the request data by combining size and payload, eg:
// '\0  \0  \0  \0 002   {   }'
// #<-->------------------------ Compression boolean (1 byte)
// #    <------------>---------- Payload size (4 bytes)
// #                    <---->-- JSON payload
func createRequestData(payload []byte) ([]byte, error) {
	size := uint32(len(payload)) // Get the size of the payload

	data := append([]byte{0x00}, make([]byte, 4)...)
	binary.BigEndian.PutUint32(data[1:], size) // Write size into the byte slice
	data = append(data, payload...)
	return data, nil
}
