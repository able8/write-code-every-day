package main

import (
	"fmt"

	"github.com/samber/lo"
)

func main() {
	list := []int{1, 2, 2, 1}

	result := lo.Uniq(list)

	fmt.Printf("%v", result)
}
