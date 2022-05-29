package typestructinterface

import "github.com/golang-jwt/jwt/v4"

// 创建一个jwt使用的密钥
var JwtKey = []byte("my_react_token_key")

type ChatMessageContent struct {
	TimeStamp          string   `json:"timeStamp,omitempty"`
	Sender             string   `json:"sender,omitempty"`
	MessageRecipientId []string `json:"messageRecipientId,omitempty"`
	ChatRoomId         string   `json:"chatRoomId,omitempty"`
	MessageTextContent string   `json:"messageTextContent,omitempty"`
}

//用户聊天是否发送者数据结构
type UsersChatroomStruct struct {
	Title  string `json:"title,omitempty"`
	Sender bool   `json:"sender,omitempty"`
}

//登陆消息结构体
type LoginRequestStrust struct {
	UserID  string `json:"userID,omitempty"`
	UserPwd string `json:"userPwd,omitempty"`
}

//token转化结构体
type JwtTokenResponseClaimsStruct struct {
	UserID string
	jwt.RegisteredClaims
}
