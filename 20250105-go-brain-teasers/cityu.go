package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	city := "Kraków"
	fmt.Println(utf8.RuneCountInString(city))
}
