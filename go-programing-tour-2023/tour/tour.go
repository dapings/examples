package main

import (
	"log"

	"github.com/dapings/examples/go-programing-tour-2023/tour/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute err: %v", err)
	}
}
