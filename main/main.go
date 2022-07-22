//http.go
package main

import (
	"log"
	"net/http"

	"server-by-chat/routes"
)

func main() {

	http.HandleFunc("/socket", routes.SocketHandler)

	http.HandleFunc("/refreshChatList", routes.RefreshChatListRouter)
	http.HandleFunc("/register", routes.RegisterRouter)
	http.HandleFunc("/login", routes.LoginRouter)
	http.HandleFunc("/refresh", routes.RefreshRouter)
	http.HandleFunc("/chatMessage", routes.ChatMessageRouter)
	http.HandleFunc("/tokenVerify", routes.TokenVerifyRouter)
	http.ListenAndServe(":9876", nil)
	log.Println("端口9876")

}
