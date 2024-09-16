package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// https://eli.thegreenplace.net/2022/go-and-proxy-servers-part-1-http-proxies/
// Appendix B: using a reverse proxy to implement a forward proxy

func main() {
	var addr = flag.String("addr", "127.0.0.1:8080", "proxy address")
	flag.Parse()

	http.HandleFunc("/", proxyHandler)
	log.Println("Starting proxy server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	target, err := url.Parse(r.URL.Scheme + "://" + r.URL.Host)
	if err != nil {
		log.Fatal(err)
	}

	reqb, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(reqb))

	p := httputil.NewSingleHostReverseProxy(target)
	p.ServeHTTP(w, r)
}
