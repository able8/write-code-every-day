package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

// getPayloadSize returns the size of the payload as a string in hexadecimal format.
func getPayloadSize(payload []byte) string {
	size := len(payload)

	formattedHex := fmt.Sprintf(`\x%02X\x%02X\x%02X\x%02X`,
		byte(size>>24),
		byte(size>>16),
		byte(size>>8),
		byte(size))

	// fmt.Printf("Formatted Size: %s\n", formattedHex)
	return formattedHex
}

func main() {
	// Define command-line flags
	dataPtr := flag.String("d", "", "JSON payload data")

	// Parse command-line flags
	flag.Parse()

	// Get the URL from the last positional argument
	args := flag.Args()
	if *dataPtr == "" || len(args) == 0 {
		fmt.Println("Error: The -d flag and a URL argument are required.")
		flag.Usage()
		return
	}

	url := args[len(args)-1] // The last argument is the URL

	// input := `{"id":5,"role":"ADMIN"}`
	// // Read payload from stdin
	// input, err := ioutil.ReadAll(os.Stdin)
	// if err != nil {
	// 	fmt.Printf("Error reading input: %s\n", err)
	// 	return
	// }

	// Generate curl command
	curlCmd := generateCurlCommand(*dataPtr, url)

	fmt.Println(curlCmd)
	// Execute the generated curl command
	output, err := executeCommand(curlCmd)
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
	} else {
		fmt.Println("Command output:")

		index := strings.Index(output, "< ")

		fmt.Println(output[index:])
	}

}

func generateCurlCommand(payload, url string) string {
	binaryHeader := `\x00` + getPayloadSize([]byte(payload)) + payload

	return fmt.Sprintf(`echo -en '%s' | curl -v -ss -k --http2 \
	-H "Content-Type: application/grpc+json" \
	-H "TE:trailers" \
	--data-binary @- https://%s`, binaryHeader, url)
}

// executeCommand executes the given command and returns its output.
func executeCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// $ echo -en '\x00\x00\x00\x00\x17{"id":1,"role":"ADMIN"}' | curl -ss -k --http2 \
//         -H "Content-Type: application/grpc+json" \
//         -H "TE:trailers" \
//         --data-binary @- \
//         https://localhost:10000/example.UserService/AddUser | od -bc
// 0000000 000 000 000 000 002 173 175
//          \0  \0  \0  \0 002   {   }
// 0000007
// $ echo -en '\x00\x00\x00\x00\x17{"id":2,"role":"GUEST"}' | curl -ss -k --http2 \
//         -H "Content-Type: application/grpc+json" \
//         -H "TE:trailers" \
//         --data-binary @- \
//         https://localhost:10000/example.UserService/AddUser | od -bc
// 0000000 000 000 000 000 002 173 175
//          \0  \0  \0  \0 002   {   }
// 0000007
$ echo -en '\x00\x00\x00\x00\x02{}' | curl -k --http2 \
        -H "Content-Type: application/grpc+json" \
        -H "TE:trailers" \
        --data-binary @- \
        --output - \
        https://localhost:10000/example.UserService/ListUsers
F{"id":1,"role":"ADMIN","create_date":"2018-07-21T20:18:21.961080119Z"}F{"id":2,"role":"GUEST","create_date":"2018-07-21T20:18:29.225624852Z"}

// go run main.go -d '{}' localhost:10000/example.UserService/ListUsers

// go run main.go -d '{"id":1,"role":"ADMIN"}' localhost:10000/example.UserService/AddUser
