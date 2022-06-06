package routes

import (
	"encoding/json"
	"log"
	"net/http"
	typestructinterface "server-by-chat/typeStructInterface"
	"sync"

	"github.com/gorilla/websocket"
)

var clientConnection sync.Map

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// 可以用来检查连接的来源
	// 这将允许从我们的 React 服务向这里发出请求。
	// 现在，我们可以不需要检查并运行任何连接
	CheckOrigin: func(r *http.Request) bool { return true },
}

func SocketHandler(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Print("Error during connection upgradation:\n\r", err)
		return
	}
	defer func() {
		conn.Close()
		log.Println("wesocket链接关闭\n\r")
	}()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息失败:\n\r", err)
			break
		}

		var chatMessageContent typestructinterface.ChatMessageContent
		err = json.Unmarshal(message, &chatMessageContent)
		if err != nil {
			log.Println("聊天消息反序列化失败\n\r")
		}
		log.Println("服务端收到的砝反序列化消息", chatMessageContent, "\n\r")
		_, ok := clientConnection.Load(chatMessageContent.Sender)
		if !ok {
			clientConnection.Store(chatMessageContent.Sender, conn)
		}
		defer func() {
			clientConnection.Delete(chatMessageContent.Sender)
			log.Println("用户下线删除", chatMessageContent.Sender, "\n\r")
		}()
		for _, v := range chatMessageContent.MessageRecipientId {
			clientConn, ok := clientConnection.Load(v)
			if ok {
				err = clientConn.(*websocket.Conn).WriteMessage(
					messageType,
					[]byte(chatMessageContent.Sender+"发给"+v+"消息了：\n\r"+chatMessageContent.MessageTextContent),
				)
				if err != nil {
					log.Println("发送消息失败:\n\r", err)
					// break
				}
			}
		}
	}
}
