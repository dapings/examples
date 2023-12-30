package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:3031/ws", nil)
	if err != nil {
		panic(err)
	}
	defer func(c *websocket.Conn, code websocket.StatusCode, reason string) {
		err := c.Close(code, reason)
		if err != nil {
			_ = c.Close(code, reason)
		}
	}(c, websocket.StatusInternalError, "内部错误！")

	err = wsjson.Write(ctx, c, "Hello WebSocket Server")
	if err != nil {
		panic(err)
	}

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}

	fmt.Printf("接收到服务端响应：%v\n", v)

	err = c.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		log.Fatal(err)
	}
}
