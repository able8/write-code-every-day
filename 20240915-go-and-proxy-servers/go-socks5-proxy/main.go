package main

import (
	"flag"

	"github.com/armon/go-socks5"
)

// https://eli.thegreenplace.net/2022/go-and-proxy-servers-part-3-socks-proxies/

// go run . -u myuser -P mypass

// http_proxy=socks5://myuser:mypass@localhost:1080 curl -v http://example.org
// https_proxy=socks5://myuser:mypass@localhost:1080 curl -vI https://example.org

type myCredentialStore struct {
	user     string
	password string
}

func (cs *myCredentialStore) Valid(user, password string) bool {
	return user == cs.user && password == cs.password
}

func main() {
	username := flag.String("u", "", "username for SOCKS5 proxy")
	password := flag.String("P", "", "password for SOCKS5 proxy")
	flag.Parse()

	auth := socks5.UserPassAuthenticator{
		Credentials: &myCredentialStore{user: *username, password: *password},
	}

	conf := &socks5.Config{
		AuthMethods: []socks5.Authenticator{auth},
	}

	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	if err := server.ListenAndServe("tcp", "127.0.0.1:1080"); err != nil {
		panic(err)
	}
}
