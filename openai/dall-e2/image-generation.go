package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"

	"github.com/sashabaranov/go-openai"
)

func generating() {
	client := openai.NewClient(os.Getenv("OPENAI_USR_TOKEN"))
	ctx := context.Background()

	// image as link
	reqURL := openai.ImageRequest{
		Prompt:         "Parrot on a skateboard performs a trick, cartoon style, natural light, high detail",
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	}
	respURL, err := client.CreateImage(ctx, reqURL)
	if err != nil {
		fmt.Printf("image generating err: %v\n", err)
		return
	}
	fmt.Println(respURL.Data[0].URL)

	// image as base64
	reqBase64 := openai.ImageRequest{
		Prompt:         "Portrait of a humanoid parrot in a classic costume, high detail, realistic light, unreal engine",
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}
	respBase64, err := client.CreateImage(ctx, reqBase64)
	if err != nil {
		fmt.Printf("image generating err: %v\n", err)
		return
	}
	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("base64 decode err: %v\n", err)
		return
	}
	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		fmt.Printf("png decode err: %v\n", err)
		return
	}

	f, err := os.Create("e2.png")
	if err != nil {
		fmt.Printf("png file creation err: %v\n", err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			_ = f.Close()
		}
	}(f)

	if err := png.Encode(f, imgData); err != nil {
		fmt.Printf("png encode err: %v\n", err)
		return
	}
	fmt.Println("the image saved as example.png")
}
