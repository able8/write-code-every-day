package main

import (
	"fmt"
	"strconv"

	"github.com/samber/lo"
)

func main() {
	list := []int64{1, 2, 3, 4}

	result := lo.FilterMap(list, func(nbr int64, index int) (string, bool) {
		return strconv.FormatInt(nbr*2, 10), nbr%2 == 0
	})

	fmt.Printf("%v", result)

}
