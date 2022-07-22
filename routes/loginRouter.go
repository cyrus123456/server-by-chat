package routes

import (
	"encoding/json"
	"log"
	"net/http"
	typestructinterface "server-by-chat/typeStructInterface"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// token有效时间
var expirationTime = time.Now().Add(5 * time.Minute)

// 临时数据库------------------------------------------------------------------------------
var usersDb = map[string]string{ //用户密码数据库
	"123": "qwe",
	"456": "qwe",
	"789": "qwe",
}

func LoginRouter(res http.ResponseWriter, req *http.Request) {
	log.Println("登陆接口")

	// 响应解码
	var loginRequestStrust typestructinterface.LoginRequestStrust
	if err := json.NewDecoder(req.Body).Decode(&loginRequestStrust); err != nil {
		log.Println("登陆接口入参对象结构解析失败👺", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("登陆接口入参对象结构%v", loginRequestStrust)

	// 获取密码
	expectedPassword, ok := usersDb[loginRequestStrust.UserID]
	if !ok {
		log.Println("没有此用户")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if loginRequestStrust.UserPwd != expectedPassword {
		log.Println("密码不正确")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("密码正确")

	jwtTokenResponseClaimsStruct := &typestructinterface.JwtTokenResponseClaimsStruct{
		UserID: loginRequestStrust.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间5分钟
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),     // 生效时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtTokenResponseClaimsStruct)
	tokenString, err := token.SignedString(typestructinterface.JwtKey)
	if err != nil {
		log.Println("token创建出错", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("token创建成功", tokenString)
	http.SetCookie(res, &http.Cookie{
		Name:  "reactToken",
		Value: tokenString,

		Path:    "/",
		Expires: expirationTime,

		Secure:   false,
		HttpOnly: false,
	})
}
