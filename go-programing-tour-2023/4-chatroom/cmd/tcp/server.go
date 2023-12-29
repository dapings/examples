package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	// 新用户到来，使用此 chan 进行登记。
	enteringChan = make(chan *User)
	// 用户离开，使用此 chan 进行登记。
	leavingChan = make(chan *User)
	// 广播专用的用户普通消息 chan，缓冲是尽可能避免出现异常情况堵塞。
	messageChan = make(chan Message, 8)
)

func main() {
	// HTTP 底层是基于TCP实现的。
	// tcp 是面向连接的协议，包括连接时的三次握手、断开时的四次挥手、传输数据时与对端进行确认接受状态的ACK、拥塞控制、失败重传等功能。
	listener, err := net.Listen("tcp", ":3030")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)

			continue
		}

		go handleConn(conn)
	}
}

// 用于记录聊天室用户，并进行消息广播：
// 1. 新用户进来；2. 用户普通消息；3. 用户离开
func broadcaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChan:
			users[user] = struct{}{}
		case user := <-leavingChan:
			delete(users, user)

			// 避免不关闭chan导致的G泄漏
			close(user.MessageChan)
		case msg := <-messageChan:
			// 给所有在线用户发送消息
			for user := range users {
				if user.ID == msg.OwnerID {
					continue
				}

				user.MessageChan <- msg.Content
			}
		}
	}
}

func handleConn(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			_ = conn.Close()
		}
	}(conn)

}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintf(conn, msg)
		if err != nil {
			log.Println("fmt.Fprintf failed", err)
		}
	}
}

type (
	User struct {
		ID          int
		Addr        string
		EnterAt     time.Time
		MessageChan chan string
	}

	// Message 给用户发送的消息。
	Message struct {
		OwnerID int
		Content string
	}
)

func (u *User) String() string {
	return u.Addr + ", UID:" + strconv.Itoa(u.ID) + ", Enter At:" +
		u.EnterAt.Format(time.DateTime+"+8000")
}

// 生产用户ID
var (
	globalID int
	idLocker sync.Mutex
)

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()

	globalID++
	return globalID
}
