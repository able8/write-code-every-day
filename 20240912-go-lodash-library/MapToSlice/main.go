// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"

	"github.com/samber/lo"
)

func main() {
	kv := map[int]int64{1: 1, 2: 2, 3: 3, 4: 4}

	result := lo.MapToSlice(kv, func(k int, v int64) string {
		return fmt.Sprintf("%d_%d", k, v)
	})

	fmt.Printf("%v", result)
}
