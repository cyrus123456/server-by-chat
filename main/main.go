//http.go
package main

import (
	"context"
	"log"
	"net/http"

	"server-by-chat/getMessageIGrpcGateway"
	"server-by-chat/routes"

	"google.golang.org/grpc"
)

func main() {

	// 连接grpc用户消息ID消息生成器服务器
	conn, err := grpc.Dial("127.0.0.1:8028", grpc.WithInsecure())
	if err != nil {
		log.Println("grpc  8028连接服务端失败: \n\r", err)
		return
	} else {
		log.Println("grpc  8028连接服务端成功:\n\r")
	}
	defer conn.Close()
	getMessageIdServiceClient := getMessageIGrpcGateway.NewGetMessageIdServiceClient(conn)
	getMessageResponse, err := getMessageIdServiceClient.GetMessageId(context.Background(), &getMessageIGrpcGateway.GetMessageRequest{UserId: "123"})
	if err != nil {
		log.Println("Grpc GetMessageId调用失败: \n\r", err)
	} else {
		log.Println("Grpc GetMessageId调用成功 \n\r", getMessageResponse)
	}

	http.HandleFunc("/socket", routes.SocketHandler)

	http.HandleFunc("/refreshChatList", routes.RefreshChatListRouter)
	http.HandleFunc("/register", routes.RegisterRouter)
	http.HandleFunc("/login", routes.LoginRouter)
	http.HandleFunc("/refresh", routes.RefreshRouter)
	http.HandleFunc("/chatMessage", routes.ChatMessageRouter)
	http.HandleFunc("/tokenVerify", routes.TokenVerifyRouter)
	http.ListenAndServe(":9876", nil)
	log.Println("端口9876\n\r")

}
