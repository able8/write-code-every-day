package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

// https://proxiesapi.com/articles/how-to-build-a-super-simple-http-proxy-in-go-in-just-20-lines-of-code

// http_proxy=http://localhost:8080 curl http://baidu.com -v
// * Uses proxy env variable http_proxy == 'http://localhost:8080'
// *   Trying [::1]:8080...
// * Connected to localhost (::1) port 8080
// > GET http://baidu.com/ HTTP/1.1

func HandleProxy(w http.ResponseWriter, r *http.Request) {
	dest_url, _ := url.Parse(r.URL.String())
	res, _ := http.Get(dest_url.String())
	io.Copy(w, res.Body)
}

func main() {
	http.HandleFunc("/", HandleProxy)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
