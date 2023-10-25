package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

func captions() {
	client := openai.NewClient(os.Getenv("OPENAI_USR_TOKEN"))
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: os.Args[1],
		Format:   openai.AudioResponseFormatSRT,
	}

	resp, err := client.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("audio captions transcription err: %v\n", err)
		return
	}
	f, err := os.Create(os.Args[1] + ".srt")
	if err != nil {
		fmt.Printf("could not open file: %v\n", err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			_ = f.Close()
		}
	}(f)

	if _, err := f.WriteString(resp.Text); err != nil {
		fmt.Printf("writing to file err: %v\n", err)
		return
	}
}
