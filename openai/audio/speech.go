package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

func speech2text() {
	client := openai.NewClient(os.Getenv("OPENAI_USR_TOKEN"))
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: "recording.mp4",
	}
	resp, err := client.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("audio transcription err: %v\n", err)
		return
	}
	fmt.Println(resp.Text)
}
