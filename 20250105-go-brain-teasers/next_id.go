package main

import (
	"fmt"
)

var (
	id = nextID()
)

func nextID() int {
	id++
	return id
}

func main() {
	fmt.Println(id)
}
