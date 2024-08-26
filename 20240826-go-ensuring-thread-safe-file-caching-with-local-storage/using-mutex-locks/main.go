package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const CACHE_FILE_PATH_TYPE_EMAIL = "/tmp/%s"

type FileWriter struct {
	filePath string
	mu       sync.Mutex
}

func NewFileWriter(filePath string) *FileWriter {
	return &FileWriter{filePath: filePath}
}

func (fw *FileWriter) Write(reader io.ReadCloser) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	file, err := os.OpenFile(fw.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		fmt.Println("Error while writing to file", err)
		return err
	}

	return nil
}

func main() {
	filename := "index.html"

	w := NewFileWriter(filename)

	var wg sync.WaitGroup

	for i := 0; i <= 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			downloadFile(filename, w)
			content := readFileContent(fmt.Sprintf(CACHE_FILE_PATH_TYPE_EMAIL, filename))
			fmt.Println(content) // Print file content
		}()
	}

	wg.Wait()
}

// Downlaod file and store locally
func downloadFile(filename string, w *FileWriter) {
	outputFilePath := fmt.Sprintf(CACHE_FILE_PATH_TYPE_EMAIL, filename)

	// Check if cache file already exist
	if fileExists(outputFilePath) {
		fmt.Println("Cache file exists")
		return
	}

	// Download the file
	fileURL := fmt.Sprintf("https://gobyexample.com/%s", filename)
	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Println("Error while downloading", err)
		return
	}
	defer resp.Body.Close()

	// Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Error creating the file", err)
		return
	}
	defer outputFile.Close()

	w.Write(resp.Body)
}

// // Copy reader buffer to file
// func write(filePath string, reader io.ReadCloser) error {
// 	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	_, err = io.Copy(file, reader)
// 	if err != nil {
// 		fmt.Println("Error while writing to file", err)
// 		return err
// 	}

// 	return nil
// }

// Read content of file
func readFileContent(filePath string) string {
	// Read the contents of the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error while reading file", err)
		return ""
	}
	return string(data)
}

// Check if file already exists
func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
