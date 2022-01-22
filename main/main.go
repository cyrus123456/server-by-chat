//http.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

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
	http.HandleFunc("/register", registerRouter)
	http.HandleFunc("/login", loginRouter)
	http.HandleFunc("/refresh", refreshRouter)
	http.HandleFunc("/chatMessage", chatMessageRouter)
	http.ListenAndServe(":8080", nil)
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
		Name:    "reactToken",
		Value:   tokenString,
		Expires: expirationTime,
	})
	// _, err = res.Write([]byte(tokenString))
	// if err != nil {
	// 	fmt.Println("token返回客户端失败", err)
	// 	res.WriteHeader(http.StatusInternalServerError)
	// }
}

func chatMessageRouter(res http.ResponseWriter, req *http.Request) {
	fmt.Println(" 聊天消息--原生HTTP路由升级webSocket")
}

func refreshRouter(res http.ResponseWriter, req *http.Request) {
	fmt.Println(" 刷新token")
}
