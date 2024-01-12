package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dapings/examples/go-programing-tour-2023/chatroom/global"
	"github.com/dapings/examples/go-programing-tour-2023/chatroom/server"
)

var (
	addr   = ":3032"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    |
   |    |    | /----\   |
   |____|    |/      \  |

 —— ChatRoom，start on：%s

`
)

func main() {
	fmt.Printf(banner, addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}

func init() {
	global.Init()
}
