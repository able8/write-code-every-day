package main

import (
	"fmt"
)

const (
	Read = 1 << iota
	Write
	Execute
)

func main() {
	mask := Read | Execute

	if mask&Execute == 0 {
		fmt.Println("can't execute")
	} else {
		fmt.Println("can execute") // will be printed
	}

	if mask&Write == 0 {
		fmt.Println("can't write") // will be printed
	} else {
		fmt.Println("can write")
	}
}
