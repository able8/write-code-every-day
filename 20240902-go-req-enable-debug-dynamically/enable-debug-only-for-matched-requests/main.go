package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
)

var debugKeyword string

var globalClient = req.EnableDumpEachRequest().OnAfterResponse(func(c *req.Client, resp *req.Response) error {
	if !resp.IsSuccessState() { // Dump and record unexpected response as error.
		return fmt.Errorf("bad response from %s %s, dump content:\n%s", resp.Request.Method, resp.Request.RawURL, resp.Dump())
	}

	// Conditional Debugging: Only display the dump content whose URL contains the specified keyword.
	if debugKeyword != "" && strings.Contains(resp.Request.RawURL, debugKeyword) {
		fmt.Println(resp.Request.Method, resp.Request.RawURL)
		fmt.Println(resp.Dump())
	}
	return nil
})

func main() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.GET("/debug", func(c *gin.Context) {
			debugKeyword = c.Query("keyword") // Enable conditional debug via API: /debug?keyword=keyword
			if debugKeyword == "off" {
				debugKeyword = ""
				fmt.Println("Debug is disabled")
			} else if debugKeyword != "" {
				fmt.Println("Start to debug", debugKeyword)
			}
		})
		r.Run("0.0.0.0:80")
	}()

	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/robots.txt",
		"https://httpbin.org/uuid",
		"https://httpbin.org/json",
		"https://httpbin.org/html",
		"https://httpbin.org/xml",
		"https://httpbin.org/headers",
		"https://httpbin.org/ip",
		"https://httpbin.org/user-agent",
		"https://api.github.com/users/imroc",
	}
	for {
		time.Sleep(5 * time.Second)
		for _, url := range urls {
			_, err := globalClient.R().
				Get(url)
			if err != nil {
				fmt.Println("ERROR:", err.Error())
				continue
			}
		}
		fmt.Println("all requests completed")
	}
}
