package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprintf(w, "HTTP, Hello")
		if err != nil {
			log.Fatalf("fmt.Fprintf failed: %v", err)
		}
	})

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println(err)

			return
		}

		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
				_ = conn.Close()
			}
		}(conn)

		// 就做一次读写
		var v any
		err = conn.ReadJSON(&v)
		if err != nil {
			log.Println(err)

			return
		}

		log.Printf("接收到客户端：%v\n", v)

		if err := conn.WriteJSON("Hello WebSocket Client"); err != nil {
			log.Println(err)

			return
		}
	})

	log.Fatal(http.ListenAndServe(":3031", nil))
}
