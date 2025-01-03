package main

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
)

type Contact struct {
	Phone string
	Email string
}

type User struct {
	Name    string
	Age     int
	contact Contact
}

func allowUnExportedInType(t reflect.Type) bool {
	if t.Name() == "User" {
		return true
	}

	return false
}

func main() {
	c1 := Contact{Phone: "123456789", Email: "dj@example.com"}
	c2 := Contact{Phone: "123456789", Email: "dj@example.com"}

	u1 := User{"dj", 18, c1}
	u2 := User{"dj", 18, c2}

	fmt.Println("u1 equals u2?", cmp.Equal(u1, u2, cmp.Exporter(allowUnExportedInType)))
}
