// sessions.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// $ go run sessions.go

// $ curl -s http://localhost:8080/secret
// Forbidden

// $ curl -s -I http://localhost:8080/login
// Set-Cookie: cookie-name=MTQ4NzE5Mz...

// $ curl -s --cookie "cookie-name=MTQ4NzE5Mz..." http://localhost:8080/secret
// The cake is a lie!

//  why logout not work in curl

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key-1")
	store = sessions.NewCookieStore(key)
)

func main() {
	http.HandleFunc("/secret", secret)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.ListenAndServe(":8080", nil)
}

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	log.Printf("secret session: %#v", session)

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	log.Printf(": %#v", session)

	// Revoke users authentication
	session.Values["authenticated"] = false
	log.Printf("after: %#v", session)
	session.Save(r, w)
}
