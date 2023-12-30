package main

import (
	"bufio"
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
				// 过滤自己发送的消息
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

	// 1. 新用户进来，构建用户的实例
	user := &User{
		ID:          GenUserID(),
		Addr:        conn.RemoteAddr().String(),
		EnterAt:     time.Now().Local(),
		MessageChan: make(chan string, 8),
	}
	// 2. 当前在一个新的G中，用来进行读操作，因此需要一个G用于写操作
	// 读写G之间可以通过channel进行通信
	go sendMessage(conn, user.MessageChan)

	// 3. 给当前用户发送欢迎信息；给所有用户告知新用户到来
	user.MessageChan <- "Welcome, " + user.String()
	msg := Message{
		OwnerID: user.ID,
		Content: "user:`" + strconv.Itoa(user.ID) + "` has enter",
	}
	messageChan <- msg

	// 4. 将该记录到全局的用户列表中，避免用锁
	enteringChan <- user

	// 控制超时用户踢出
	var userActive = make(chan struct{})
	go func() {
		d := 1 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				// 自动踢出不活跃用户
				_ = conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	// 5. 循环读取用户的输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = strconv.Itoa(user.ID) + ":" + input.Text()
		messageChan <- msg

		// 用户活跃
		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}

	// 6. 用户离开
	leavingChan <- user
	msg.Content = "user:`" + strconv.Itoa(user.ID) + "` has left"
	messageChan <- msg
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
