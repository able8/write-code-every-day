package main

import "github.com/armon/go-socks5"

// https://github.com/armon/go-socks5
// curl -Iv -x socks5://localhost:8888 https://example.com

func main() {
	// Create a SOCKS5 server
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", "0.0.0.0:8888"); err != nil {
		panic(err)
	}
}
