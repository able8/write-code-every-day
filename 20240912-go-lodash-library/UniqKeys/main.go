// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"

	"github.com/samber/lo"
)

func main() {
	kv1 := map[string]int{"foo": 1, "bar": 2}
	kv2 := map[string]int{"bar": 3, "baz": 4}

	result := lo.UniqKeys(kv1, kv2)

	fmt.Printf("%v", result)
}
