package logic

import (
	"expvar"
	"log"

	"github.com/dapings/examples/go-programing-tour-2023/chatroom/global"
)

func init() {
	expvar.Publish("message_queue", expvar.Func(calcMessageQueueLen))
}

func calcMessageQueueLen() any {
	return 0
}

var (
	// Broadcaster 全局只有一个广播器，使用饿汉式单例模式实现。
	Broadcaster = &broadcaster{
		users: make(map[string]*User),

		enteringChannel: make(chan *User),
		leavingChannel:  make(chan *User),
		messageChannel:  make(chan *Message, global.MessageQueueLen),

		checkUserChannel:      make(chan string),
		checkUserCanInChannel: make(chan bool),

		requestUsersChannel: make(chan struct{}),
		usersChannel:        make(chan []*User),
	}
)

// broadcaster 广播器
type broadcaster struct {
	// 所有聊天室用户
	users map[string]*User

	// 所有channel统一管理，可以避免外部乱用

	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	// 判断此昵称用户是否可进入聊天室(重复与否)：true 能，false 不能
	checkUserChannel      chan string
	checkUserCanInChannel chan bool

	// 获取用户列表
	requestUsersChannel chan struct{}
	usersChannel        chan []*User
}

// Start 启动广播器
// 需要在一个新G中运行，因为它不会返回
// 最佳实践(没有在自己内部开启新的G)：应该让调用者决定并发(启动新的G)，这样它清楚自己在干什么。
func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			// 新用户进入
			b.users[user.NickName] = user

			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			// 用户离开
			delete(b.users, user.NickName)
			// 关闭MessageChannel后，使相应的G结束，从而避免G泄漏，导致内存的泄漏
			user.CloseMessageChannel()
		case msg := <-b.messageChannel:
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}

				user.MessageChannel(msg)
			}

			OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUsersChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("broadcast message queue >= global queue len")
		// TODO: 队列满了，还要继续写数据，或等待有空位？
	}

	b.messageChannel <- msg
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname

	return <-b.checkUserCanInChannel
}

func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}

	return <-b.usersChannel
}
