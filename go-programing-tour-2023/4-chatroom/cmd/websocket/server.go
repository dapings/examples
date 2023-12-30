package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprintf(w, "HTTP, Hello")
		if err != nil {
			log.Fatalf("fmt.Fprintf failed: %v", err)
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		conn, err := websocket.Accept(w, req, nil)
		if err != nil {
			log.Println(err)

			return
		}

		defer func(conn *websocket.Conn, code websocket.StatusCode, reason string) {
			err := conn.Close(code, reason)
			if err != nil {
				_ = conn.Close(code, reason)
			}
		}(conn, websocket.StatusInternalError, "内部出错了！")

		ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			log.Println(err)

			return
		}
		log.Printf("接收到客户端：%v\n", v)

		err = wsjson.Write(ctx, conn, "Hello WebSocket Client")
		if err != nil {
			log.Println(err)

			return
		}

		err = conn.Close(websocket.StatusNormalClosure, "")
		if err != nil {
			log.Fatal(err)
		}
	})

	log.Fatal(http.ListenAndServe(":3031", nil))
}
