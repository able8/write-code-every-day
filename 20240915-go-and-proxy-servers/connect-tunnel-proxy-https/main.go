// Implements a tunneling forward proxy for CONNECT requests.
//
// This proxy serves over plain HTTP and only supports CONNECT requests to
// arbitrary targets; it does not support regular HTTP proxying.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
)

// https_proxy=http://localhost:9999 curl https://baidu.com -vI

// We see that curl contacted our proxy, sending it a CONNECT request for https://example.org.
// Our proxy replied with 200 OK and then curl proceeded to perform a TLS handshake with the destination server,
// through the proxy. It works!

func main() {
	var addr = flag.String("addr", "127.0.0.1:9999", "proxy address")
	flag.Parse()

	proxy := &forwardProxy{}

	log.Println("Starting proxy server on", *addr)
	if err := http.ListenAndServe(*addr, proxy); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type forwardProxy struct {
}

func (p *forwardProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		proxyConnect(w, req)
	} else {
		http.Error(w, "this proxy only supports CONNECT", http.StatusMethodNotAllowed)
	}
}

// This code reads the destination address from the request and establishes a new TCP connection.
//	It then... hijacks the client connection?
//

func proxyConnect(w http.ResponseWriter, req *http.Request) {
	log.Printf("CONNECT requested to %v (from %v)", req.Host, req.RemoteAddr)
	targetConn, err := net.Dial("tcp", req.Host)
	if err != nil {
		log.Println("failed to dial to target", req.Host)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Fatal("http server doesn't support hijacking connection")
	}

	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Fatal("http hijacking failed")
	}

	// In this case of tunneling traffic in the proxy, the raw TCP connection is exactly what we need.
	// So we end up with two TCP connections - one with the client and one with the destination server. The next step is to hook them up together.
	// We start two goroutines - one for each direction; tunnelConn does this:

	log.Println("tunnel established")
	go tunnelConn(targetConn, clientConn)
	go tunnelConn(clientConn, targetConn)
}

func tunnelConn(dst io.WriteCloser, src io.ReadCloser) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
}
