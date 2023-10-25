package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func supportCtx() {
	client := openai.NewClient(os.Getenv("OPENAI_USR_TOKEN"))
	msgs := make([]openai.ChatCompletionMessage, 0)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("conversation")
	fmt.Println("======================")

	for {
		fmt.Print("-> ")
		txt, _ := reader.ReadString('\n')
		// convert CRLF to LF
		txt = strings.Replace(txt, "\n", "", -1)
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: txt,
		})

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: msgs,
			},
		)
		if err != nil {
			fmt.Printf("chatgpt completion err: %v\n", err)
			continue
		}

		content := resp.Choices[0].Message.Content
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
		fmt.Println(content)
	}
}
