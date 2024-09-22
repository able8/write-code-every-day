package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=18"`
}

func main() {
	user := &User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Age:   17,
	}

	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		fmt.Println("Validation failed:")
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Printf("Field: %s, Error: %s\n", e.Field(), e.Tag())
		}
	} else {
		fmt.Println("Validation succeeded")
	}
}
