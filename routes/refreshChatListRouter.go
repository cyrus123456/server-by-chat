package routes

import (
	"encoding/json"
	"log"
	"net/http"
	typestructinterface "server-by-chat/typeStructInterface"
	"time"
)

var usersChatroomDb = map[string][]typestructinterface.UsersChatroomStruct{ //用户聊天室列表数据库
	"123": {
		{
			Title:  "123_456_789",
			Sender: true,
		},
		{
			Title:  "123_456",
			Sender: false,
		},
	},
	"456": {
		{
			Title:  "123_456_789",
			Sender: true,
		},
		{
			Title:  "123_456",
			Sender: false,
		},
	},
	"789": {
		{
			Title:  "123_456_789",
			Sender: true,
		},
	},
}

var timeNowFormat = time.Now().Format("2006-01-02 15:04:05") //当前时间

var chatroomDb = map[string][]typestructinterface.ChatMessageContent{ //聊天室列表数据库
	"123_456_789": {
		{
			TimeStamp:          timeNowFormat,
			Sender:             "123",
			MessageRecipientId: make([]string, 0),
			ChatRoomId:         "123_456_789",
			MessageTextContent: "大家好",
		},
	},
	"123_456": {
		{
			TimeStamp:          timeNowFormat,
			Sender:             "123",
			MessageRecipientId: make([]string, 0),
			ChatRoomId:         "123_456",
			MessageTextContent: "你好",
		},
		{
			TimeStamp:          timeNowFormat,
			Sender:             "456",
			MessageRecipientId: make([]string, 0),
			ChatRoomId:         "123_456",
			MessageTextContent: "enen,你好",
		},
	},
}

func RefreshChatListRouter(res http.ResponseWriter, req *http.Request) {

	type UidStruct struct {
		Uid string `json:"uid,omitempty"`
	}
	uidStruct := UidStruct{}
	if err := json.NewDecoder(req.Body).Decode(&uidStruct); err != nil {
		log.Println("初始化聊天信息接口入参对象结构解析失败👺", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("初始化聊天信息接口入参对象", uidStruct.Uid)

	type ResStruct struct {
		UsersChatroomDb []typestructinterface.UsersChatroomStruct
		ChatroomDb      map[string][]typestructinterface.ChatMessageContent
	}
	resStruct := ResStruct{
		UsersChatroomDb: usersChatroomDb[uidStruct.Uid],
		ChatroomDb:      chatroomDb,
	}
	resByyte, err := json.Marshal(resStruct)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Write(resByyte)
}
