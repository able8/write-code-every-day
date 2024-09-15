package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// https://eli.thegreenplace.net/2022/go-and-proxy-servers-part-1-http-proxies/

// Setting up our own reverse proxy in Go
// The Go standard library makes it pretty easy with the net/http/httputil package
// which provides a ReverseProxy type and a simple constructor named NewSingleHostReverseProxy
// to implement the kind of forwarding shown in the previous section.
// To create a simple reverse proxy that forwards from one address to another, all we need to write is:

// And with our debugging server still listening on 8080, curl requests will be redirected by the proxy to the server.
func main() {
	fromAddr := flag.String("from", "127.0.0.1:9090", "proxy's listening address")
	toAddr := flag.String("to", "127.0.0.1:8080", "the address this proxy will forward to")
	flag.Parse()

	toUrl := parseToUrl(*toAddr)
	proxy := httputil.NewSingleHostReverseProxy(toUrl)
	log.Println("Starting proxy server on", *fromAddr)
	if err := http.ListenAndServe(*fromAddr, proxy); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// parseToUrl parses a "to" address to url.URL value
func parseToUrl(addr string) *url.URL {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}
	toUrl, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	return toUrl
}
