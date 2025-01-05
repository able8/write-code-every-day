package main

import (
	"fmt"
)

type FilePerm uint16 // 16 flags are enough

const (
	Read FilePerm = 1 << iota
	Write
	Execute
)

// String implements fmt.Stringer interface
func (p FilePerm) String() string {
	switch p {
	case Read:
		return "read"
	case Write:
		return "write"
	case Execute:
		return "execute"
	}

	return fmt.Sprintf("unknown FilePerm: %d", p) // don't use %s here :)
}

func main() {
	fmt.Println(Execute)        // execute
	fmt.Printf("%d\n", Execute) // 4
}
