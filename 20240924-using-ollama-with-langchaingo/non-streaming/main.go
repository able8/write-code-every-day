package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// ollama run llama3.1
// ollama serve

func main() {
	llm, err := ollama.New(ollama.WithModel("llama3.1"))
	if err != nil {
		log.Fatal(err)
	}

	query := "very briefly, tell me the difference between fall and autumn"

	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:\n", completion)
}
