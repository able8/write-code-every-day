package main

import (
	"fmt"
)

func main() {
	r1 := '™'
	fmt.Println(r1)        // 8482
	fmt.Printf("%c\n", r1) // ™

	r2 := '\x61'           // \x - 2 digits
	fmt.Printf("%c\n", r2) // a

	r3 := '\u2122'         // \u - 4 digits (8482 in hex)
	fmt.Printf("%c\n", r3) // ™

	r4 := '\U00002122'     // \U - 8 digits
	fmt.Printf("%c\n", r4) // ™
}
