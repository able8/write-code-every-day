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

// implement a very simplistic "load-balancing" reverse proxy
// that round-robins between two backend servers.
// Here's the interesting code snippet (the full code is here):

func main() {
	fromAddr := flag.String("from", "127.0.0.1:9090", "proxy's listening address")
	toAddr1 := flag.String("to1", "127.0.0.1:8080", "the first address this proxy will forward to")
	toAddr2 := flag.String("to2", "127.0.0.1:8081", "the second address this proxy will forward to")
	flag.Parse()

	toUrl1 := parseToUrl(*toAddr1)
	toUrl2 := parseToUrl(*toAddr2)

	proxy := loadBalancingReverseProxy(toUrl1, toUrl2)
	log.Println("Starting proxy server on", *fromAddr)
	if err := http.ListenAndServe(*fromAddr, proxy); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func loadBalancingReverseProxy(target1, target2 *url.URL) *httputil.ReverseProxy {
	var targetNum = 1

	director := func(req *http.Request) {
		var target *url.URL
		// Simple round robin between the two targets
		if targetNum == 1 {
			target = target1
			targetNum = 2
		} else {
			target = target2
			targetNum = 1
		}

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
		// For simplicity, we don't handle RawQuery or the User-Agent header here:
		// see the full code of NewSingleHostReverseProxy for an example of doing
		// that.
	}
	return &httputil.ReverseProxy{Director: director}
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
