package main

import (
	"fmt"
	"net/http"
	"sync"
)

func fetchURL(url string, ch chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Error fetching %s: %v", url, err)
		return
	}
	ch <- fmt.Sprintf("Fetched %s with status %s", url, resp.Status)
}

func main() {
	urls := []string{
		"http://example.com",
		"http://example.org",
		"http://example.net",
	}

	var wg sync.WaitGroup
	ch := make(chan string)

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fetchURL(url, ch)
		}()
	}

	go func() {
		wg.Wait()
		close(ch) // 关闭 channel
	}()

	for msg := range ch {
		fmt.Println(msg) // 打印每个结果
	}
}
