package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
)

func streaming() {
	client := openai.NewClient(os.Getenv("OPENAI_USR_TOKEN"))
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     openai.GPT3Ada,
		Prompt:    "heap sort py code",
		MaxTokens: 5,
		Stream:    true,
	}
	stream, err := client.CreateCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("gpt3 streaming err: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("streaming finished")
			return
		}
		if err != nil {
			fmt.Printf("streaming err: %v\n", err)
			return
		}
		fmt.Printf("streaming resp: %v\n", resp)
	}
}
