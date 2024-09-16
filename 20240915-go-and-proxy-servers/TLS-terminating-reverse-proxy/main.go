package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// https://eli.thegreenplace.net/2022/go-and-proxy-servers-part-2-https-proxies/

// Why won't HTTP proxies "just work" for HTTPS?
// The reason is that an HTTPS client expects to talk to a specific server,
// and will look for a valid certificate from that server to start sending information.

func main() {
	fromAddr := flag.String("from", "127.0.0.1:9090", "proxy's listening address")
	toAddr := flag.String("to", "127.0.0.1:8080", "the address this proxy will forward to")
	certFile := flag.String("certfile", "cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "key.pem", "key PEM file")
	flag.Parse()

	toUrl := parseToUrl(*toAddr)
	proxy := httputil.NewSingleHostReverseProxy(toUrl)

	srv := &http.Server{
		Addr:    *fromAddr,
		Handler: proxy,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}
	log.Println("Starting proxy server on", *fromAddr)
	if err := srv.ListenAndServeTLS(*certFile, *keyFile); err != nil {
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
