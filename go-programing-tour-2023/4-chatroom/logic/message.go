package logic

import (
	"time"
)

// Message 给用户发送的消息。
type Message struct {
	// 那个用户发送的消息
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	ClientSentTime time.Time `json:"client_sent_time"`

	// 消息@了谁
	Ats []string `json:"ats"`

	// 用户列表不通过 websocket 下发
	// Users []*User `json:"users"`
}
