//http.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// 可以用来检查连接的来源
	// 这将允许从我们的 React 服务向这里发出请求。
	// 现在，我们可以不需要检查并运行任何连接
	CheckOrigin: func(r *http.Request) bool { return true },
}

// 创建一个jwt使用的密钥
var jwtKey = []byte("my_react_token_key")

// 临时数据库------------------------------------------------------------------------------
var usersDb = map[string]string{ //用户密码数据库
	"123": "qwe",
	"456": "qwe",
	"789": "qwe",
}

type usersChatroomStruct struct {
	Title  string `json:"title,omitempty"`
	Sender bool   `json:"sender,omitempty"`
}

var usersChatroomDb = map[string][]usersChatroomStruct{ //用户聊天室列表数据库
	"123": []usersChatroomStruct{
		usersChatroomStruct{
			Title:  "123_456_789",
			Sender: true,
		},
		usersChatroomStruct{
			Title:  "123_456",
			Sender: false,
		},
	},
	"456": []usersChatroomStruct{
		usersChatroomStruct{
			Title:  "123_456_789",
			Sender: true,
		},
		usersChatroomStruct{
			Title:  "123_456",
			Sender: false,
		},
	},
	"789": []usersChatroomStruct{
		usersChatroomStruct{
			Title:  "123_456_789",
			Sender: true,
		},
	},
}

type ChatMessageContent struct {
	TimeStamp          string   `json:"timeStamp,omitempty"`
	Sender             string   `json:"sender,omitempty"`
	MessageRecipientId []string `json:"messageRecipientId,omitempty"`
	ChatRoomId         string   `json:"chatRoomId,omitempty"`
	MessageTextContent string   `json:"messageTextContent,omitempty"`
}

var timeNowFormat = time.Now().Format("2006-01-02 15:04:05") //当前时间
var chatroomDb = map[string][]ChatMessageContent{            //聊天室列表数据库
	"123_456_789": []ChatMessageContent{
		ChatMessageContent{
			TimeStamp:          timeNowFormat,
			Sender:             "123",
			MessageRecipientId: make([]string, 0),
			ChatRoomId:         "123_456_789",
			MessageTextContent: "大家好",
		},
	},
	"123_456": []ChatMessageContent{
		ChatMessageContent{
			TimeStamp:          timeNowFormat,
			Sender:             "123",
			MessageRecipientId: make([]string, 0),
			ChatRoomId:         "123_456",
			MessageTextContent: "你好",
		},
		ChatMessageContent{
			TimeStamp:          timeNowFormat,
			Sender:             "456",
			MessageRecipientId: make([]string, 0),
			ChatRoomId:         "123_456",
			MessageTextContent: "enen,你好",
		},
	},
}

type LoginRequestStrust struct {
	UserID  string `json:"userID,omitempty"`
	UserPwd string `json:"userPwd,omitempty"`
}

//token转化结构体
type JwtTokenResponseClaimsStruct struct {
	UserID string
	jwt.StandardClaims
}

var clientConnection sync.Map

func main() {
	// 练习
	// practiceInterview.Test()

	log.Println("端口9876\n\r")

	http.HandleFunc("/socket", socketHandler)

	http.HandleFunc("/refreshChatList", refreshChatListRouter)
	http.HandleFunc("/register", registerRouter)
	http.HandleFunc("/login", loginRouter)
	http.HandleFunc("/refresh", refreshRouter)
	http.HandleFunc("/chatMessage", chatMessageRouter)
	http.HandleFunc("/tokenVerify", tokenVerifyRouter)
	http.ListenAndServe(":9876", nil)
}

func socketHandler(res http.ResponseWriter, req *http.Request) {
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
		var chatMessageContent ChatMessageContent
		err = json.Unmarshal(message, &chatMessageContent)
		if err != nil {
			log.Println("聊天消息反序列化失败\n\r")
		}
		log.Println("服务端收到的砝反序列化消息", chatMessageContent, "\n\r")
		_, ok := clientConnection.Load(chatMessageContent.Sender)
		if !ok {
			clientConnection.Delete(chatMessageContent.Sender)
			clientConnection.Store(chatMessageContent.Sender, conn)
		}
		defer func() {
			log.Println("用户下线", chatMessageContent.Sender, "\n\r")
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

func refreshChatListRouter(res http.ResponseWriter, req *http.Request) {

	type UidStruct struct {
		Uid string `json:"uid,omitempty"`
	}
	uidStruct := UidStruct{}
	if err := json.NewDecoder(req.Body).Decode(&uidStruct); err != nil {
		log.Println("初始化聊天信息接口入参对象结构解析失败\n\r", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("初始化聊天信息接口入参对象\r\n", uidStruct.Uid)

	type ResStruct struct {
		UsersChatroomDb []usersChatroomStruct
		ChatroomDb      map[string][]ChatMessageContent
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

func registerRouter(res http.ResponseWriter, req *http.Request) {
	log.Println("注册--原生HTTP路由")
}

func loginRouter(res http.ResponseWriter, req *http.Request) {
	log.Println("登陆接口\n\r")

	// 响应解码
	var loginRequestStrust LoginRequestStrust
	if err := json.NewDecoder(req.Body).Decode(&loginRequestStrust); err != nil {
		log.Println("登陆接口入参对象结构解析失败\n\r", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("登陆接口入参对象结构%v\n", loginRequestStrust)

	// 获取密码
	expectedPassword, ok := usersDb[loginRequestStrust.UserID]
	if !ok {
		log.Println("没有此用户\n\r")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if loginRequestStrust.UserPwd != expectedPassword {
		log.Println("密码不正确\n\r")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("密码正确\n\r")

	expirationTime := time.Now().Add(5 * time.Minute)
	jwtTokenResponseClaimsStruct := &JwtTokenResponseClaimsStruct{
		UserID: loginRequestStrust.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtTokenResponseClaimsStruct)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("token创建出错\n\r", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("token创建成功\n\r", tokenString)
	http.SetCookie(res, &http.Cookie{
		Name:  "reactToken",
		Value: tokenString,

		Path:    "/",
		Expires: expirationTime,

		Secure:   false,
		HttpOnly: false,
	})
}

func refreshRouter(res http.ResponseWriter, req *http.Request) {
	log.Println(" 刷新token\n\r")
	httpCookie, err := req.Cookie("reactToken")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("未设置cookie\n\r")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Println("请求错误\n\r")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("成功获取cookie\n\r", httpCookie.Value)
	jwtTokenResponseClaimsStruct := &JwtTokenResponseClaimsStruct{}
	okToken, err := jwt.ParseWithClaims(httpCookie.Value, jwtTokenResponseClaimsStruct, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	log.Println("验证token\n\r", okToken)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		res.WriteHeader(http.StatusBadGateway)
		return
	}
	if !okToken.Valid {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	if time.Unix(jwtTokenResponseClaimsStruct.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		res.WriteHeader(http.StatusBadRequest)
		log.Println("jwtToken到期还有很长时间\n\r")
		return
	}
	log.Println("开始颁发新的jwttoken\n\r")
	expirationTime := time.Now().Add(5 * time.Minute)
	log.Println("逾期时间对象\n\r", expirationTime)
	jwtTokenResponseClaimsStruct.StandardClaims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtTokenResponseClaimsStruct)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("token创建出错\n\r", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("token创建成功\n\r", tokenString)
	http.SetCookie(res, &http.Cookie{
		Name:  "reactToken",
		Value: tokenString,

		Path:    "/",
		Expires: expirationTime,

		Secure:   true,
		HttpOnly: true,
	})
}

func chatMessageRouter(res http.ResponseWriter, req *http.Request) {
	log.Println(" 聊天消息--原生HTTP路由升级webSocket\n\r")
}

func tokenVerifyRouter(res http.ResponseWriter, req *http.Request) {
	log.Println(" token验证\n\r")
	httpCookie, err := req.Cookie("reactToken")
	if err != nil {
		if err == http.ErrNoCookie {
			// 如果未设置cookie，则返回未授权状态
			log.Println("未设置Cookie\n\r")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		// 对于其他类型的错误，返回错误的请求状态。
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("成功获取token\n\r", httpCookie.Value)
	jwtTokenResponseClaimsStruct := &JwtTokenResponseClaimsStruct{}
	okToken, err := jwt.ParseWithClaims(httpCookie.Value, jwtTokenResponseClaimsStruct, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	log.Println("验证token\n\r", okToken)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if !okToken.Valid {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 最后，将欢迎消息以及令牌中的用户名返回给用户
	res.Write([]byte(fmt.Sprintf("Welcome %s!\n\r", jwtTokenResponseClaimsStruct.UserID)))
}
