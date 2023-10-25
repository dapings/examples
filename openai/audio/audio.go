package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please provide a filename to convert to text")
		return
	}

	if _, err := os.Stat(os.Args[1]); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("file %s not exits\n", os.Args[1])
		return
	}

	client := openai.NewClient(os.Getenv("OPENAI_USR_TOKEN"))
	ctx := context.Background()
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: os.Args[1],
	}
	resp, err := client.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("voice transcription err: %s\n", err)
		return
	}
	fmt.Println(resp.Text)
}
