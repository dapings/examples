package server

import (
	"log"
	"net/http"

	"github.com/dapings/examples/go-programing-tour-2023/chatroom/logic"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept 从客户端接受 Websocket 握手，并将连接升级到 Websocket。
	// 如果Origin域与主机不同，Accept将拒绝握手，除非设置了InsecureSkipVerify选项(通过第三个参数AcceptOptions设置)。
	// 换句话说，默认情况下，它不允许跨源请求。如果发生错误，Accept将始终写入适当的响应。
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Println("websocket accept error:", err)

		return
	}

	// 1. 新用户进来，构建此用户的实例
	token := req.FormValue("token")
	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal: ", nickname)
		_ = wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度：2-20"))
		_ = conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")

		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("nickname exists： ", nickname)
		_ = wsjson.Write(req.Context(), conn, logic.NewErrorMessage("昵称已存在！"))
		_ = conn.Close(websocket.StatusUnsupportedData, "nickname exists!")

		return
	}

	userHasToken := logic.NewUser(conn, token, nickname, req.RemoteAddr)

	// 2. 开启给用户发送消息的G
	go userHasToken.SendMessage(req.Context())

	// 3. 给当前用户发送欢迎信息
	userHasToken.MessageChannel(logic.NewWelcomeMessage(userHasToken))

	// 避免token泄漏
	tmpUser := *userHasToken
	user := &tmpUser
	user.Token = ""

	// 给所有用户告知新用户到来
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	// 4. 将此用户加入广播器的用户列表中
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "join chat")

	// 5. 接收用户消息
	err = user.ReceiveMessage(req.Context())

	// 6. 用户离开
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, "leaves chat")

	// 根据读取时的错误，执行不同的Close
	if err == nil {
		_ = conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		_ = conn.Close(websocket.StatusInternalError, "read from client error")
	}
}
