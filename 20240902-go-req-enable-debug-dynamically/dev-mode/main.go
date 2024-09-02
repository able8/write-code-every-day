package main

import (
	"github.com/imroc/req/v3"
)

func main() {
	client := req.C().DevMode()
	client.R().Get("https://httpbin.org/get")

	client = req.C().EnableDebugLog()
	client.R().Get("http://baidu.com/s?wd=req")

}
