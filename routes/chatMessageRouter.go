package routes

import (
	"log"
	"net/http"
)

func ChatMessageRouter(res http.ResponseWriter, req *http.Request) {
	log.Println(" 聊天消息--原生HTTP路由升级webSocket")
}
