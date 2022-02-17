//http.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	practiceinterview "server-by-chat/practiceInterview"
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

// 临时数据库
var usersDb = map[string]string{
	"123": "qwe",
	"456": "qwe",
}

type LoginRequestStrust struct {
	UserID  string `json:"userID"`
	UserPwd string `json:"userPwd"`
}

//token转化结构体
type JwtTokenResponseClaimsStruct struct {
	UserID string
	jwt.StandardClaims
}

func main() {
	// 练习
	practiceinterview.Test()

	http.HandleFunc("/socket", socketHandler)
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
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		fmt.Println(string(message))
		log.Printf("Received: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}
}

func registerRouter(res http.ResponseWriter, req *http.Request) {
	fmt.Println("注册--原生HTTP路由")
}

func loginRouter(res http.ResponseWriter, req *http.Request) {
	fmt.Println("登陆接口")

	// 响应解码
	var loginRequestStrust LoginRequestStrust
	if err := json.NewDecoder(req.Body).Decode(&loginRequestStrust); err != nil {
		fmt.Println("登陆接口入参对象结构解析失败", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("登陆接口入参对象结构%v\n", loginRequestStrust)

	// 获取密码
	expectedPassword, ok := usersDb[loginRequestStrust.UserID]
	if !ok {
		fmt.Println("没有此用户")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if loginRequestStrust.UserPwd != expectedPassword {
		fmt.Println("密码不正确")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println("密码正确")

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
		fmt.Println("token创建出错", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("token创建成功", tokenString)
	http.SetCookie(res, &http.Cookie{
		Name:  "reactToken",
		Value: tokenString,

		Path:    "/",
		Expires: expirationTime,

		Secure:   true,
		HttpOnly: true,
	})
}

func refreshRouter(res http.ResponseWriter, req *http.Request) {
	fmt.Println(" 刷新token")
	httpCookie, err := req.Cookie("reactToken")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("未设置cookie")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Println("请求错误")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("成功获取cookie", httpCookie.Value)
	jwtTokenResponseClaimsStruct := &JwtTokenResponseClaimsStruct{}
	okToken, err := jwt.ParseWithClaims(httpCookie.Value, jwtTokenResponseClaimsStruct, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	fmt.Println("验证token", okToken)
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
		fmt.Println("jwtToken到期还有很长时间")
		return
	}
	fmt.Println("开始颁发新的jwttoken")
	expirationTime := time.Now().Add(5 * time.Minute)
	fmt.Println("逾期时间对象", expirationTime)
	jwtTokenResponseClaimsStruct.StandardClaims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtTokenResponseClaimsStruct)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("token创建出错", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("token创建成功", tokenString)
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
	fmt.Println(" 聊天消息--原生HTTP路由升级webSocket")
}

func tokenVerifyRouter(res http.ResponseWriter, req *http.Request) {
	fmt.Println(" token验证")
	httpCookie, err := req.Cookie("reactToken")
	if err != nil {
		if err == http.ErrNoCookie {
			// 如果未设置cookie，则返回未授权状态
			fmt.Println("未设置Cookie")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		// 对于其他类型的错误，返回错误的请求状态。
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("成功获取token", httpCookie.Value)
	jwtTokenResponseClaimsStruct := &JwtTokenResponseClaimsStruct{}
	okToken, err := jwt.ParseWithClaims(httpCookie.Value, jwtTokenResponseClaimsStruct, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	fmt.Println("验证token", okToken)
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
	res.Write([]byte(fmt.Sprintf("Welcome %s!", jwtTokenResponseClaimsStruct.UserID)))
}
