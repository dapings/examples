package server

import (
	"net/http"

	"github.com/dapings/examples/go-programing-tour-2023/chatroom/logic"
)

func RegisterHandle() {
	// 广播消息处理
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/user_list", userListHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}
