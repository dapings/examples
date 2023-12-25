package main

import (
	_ "net/http/pprof"
	
	"github.com/dapings/examples/go-programing-tour-2023/crawler/cmd"
)

func main() {
	cmd.Execute()
}
