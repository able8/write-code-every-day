// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"

	"github.com/samber/lo"
)

func main() {
	kv := map[string]int{"foo": 1, "bar": 2, "baz": 3}

	result := lo.Entries(kv)

	fmt.Printf("%v", result)
}
