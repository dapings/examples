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

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 20,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "ChatGPT streaming completion example",
			},
		},
		Stream: true,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream err: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("streaming resp: ")
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream err: %v\n", err)
			return
		}

		fmt.Printf(resp.Choices[0].Delta.Content)
	}
}
