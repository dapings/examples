package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

func main() {
	client := openai.NewClient(os.Getenv("OPENAI_USR_TOKEN"))
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     openai.GPT3Ada,
		MaxTokens: 5,
		Prompt:    "quick sort py code",
	}
	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Printf("GPT3 completion err: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Text)
}
