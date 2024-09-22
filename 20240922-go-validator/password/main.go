package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func main() {
	var password = "securepass"

	v := validator.New()

	// 👇
	err := v.Var(password, "required,min=8,containsany=!@#?*")
	if err != nil {
		fmt.Println(err)
	}
}
