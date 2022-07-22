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
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer func() {
		// 异常关闭链接
		conn.Close()
		log.Println("wesocket链接关闭")
	}()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息失败👺:", err)
			break
		}

		var chatMessageContent typestructinterface.ChatMessageContent
		err = json.Unmarshal(message, &chatMessageContent)
		if err != nil {
			log.Println("聊天消息反序列化失败👺")
		}
		log.Println("服务端收到的砝反序列化消息", chatMessageContent, "")
		_, ok := clientConnection.Load(chatMessageContent.Sender)
		if !ok {
			// 避免重复保存用户链接
			clientConnection.Store(chatMessageContent.Sender, conn)
		}
		defer func() {
			// 链接断开删除用户
			clientConnection.Delete(chatMessageContent.Sender)
			log.Println("用户下线删除", chatMessageContent.Sender, "")
		}()
		for _, v := range chatMessageContent.MessageRecipientId {
			clientConn, ok := clientConnection.Load(v)
			if ok {
				// 如果用户在线
				err = clientConn.(*websocket.Conn).WriteMessage(
					messageType,
					[]byte(chatMessageContent.Sender+"发给"+v+"消息了："+chatMessageContent.MessageTextContent),
				)
				if err != nil {
					log.Println("发送消息失败👺:", err)
					// break
				}
			} else {
				//如果用户不在线，漫游信息临时redis缓存，个个不同用户端上线后拉取信息

			}
		}
	}
}
