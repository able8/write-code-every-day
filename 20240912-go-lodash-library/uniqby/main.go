package main

import (
	"fmt"

	"github.com/samber/lo"
)

func main() {
	list := []int{0, 1, 2, 3, 4, 5}

	result := lo.UniqBy(list, func(i int) int {
		return i % 3
	})

	fmt.Printf("%v", result)
}
