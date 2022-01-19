//http.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 临时数据库
var usersDb = map[string]string{
	"123": "password1",
	"456": "password2",
}

type loginApiPramsStrust struct {
	UserID  string `json:"userID"`
	UserPwd string `json:"userPwd"`
}

//token转化结构体
type ClaimsStruct struct {
	UserID  string
	UserPwd string
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
	fmt.Fprint(res, " 注册--原生HTTP路由")
}

func loginRouter(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "登陆--原生HTTP路由")

	// 响应解码
	var loginApiPramsStrust loginApiPramsStrust
	if err := json.NewDecoder(req.Body).Decode(&loginApiPramsStrust); err != nil {
		fmt.Println("登陆接口入参对象结构解析失败")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	claimsStruct := ClaimsStruct{
		UserID:  loginApiPramsStrust.UserID,
		UserPwd: loginApiPramsStrust.UserPwd,
	}
	// 获取密码
	expectedPassword, ok := usersDb[claimsStruct.UserID]
	if !ok {
		fmt.Println("没有此用户")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if claimsStruct.UserPwd != expectedPassword {
		fmt.Println("密码不正确")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsStruct)
	tokenString, err := token.SigningString()
	if err != nil {
		fmt.Println("token创建出错")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(5 * time.Minute),
	})
}

func chatMessageRouter(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, " 聊天消息--原生HTTP路由升级webSocket")
}

func refreshRouter(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, " 刷新token")
}
