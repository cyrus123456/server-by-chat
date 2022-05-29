package routes

import (
	"fmt"
	"log"
	"net/http"
	typestructinterface "server-by-chat/typeStructInterface"

	"github.com/golang-jwt/jwt/v4"
)

func TokenVerifyRouter(res http.ResponseWriter, req *http.Request) {
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
	jwtTokenResponseClaimsStruct := &typestructinterface.JwtTokenResponseClaimsStruct{}
	okToken, err := jwt.ParseWithClaims(httpCookie.Value, jwtTokenResponseClaimsStruct, func(token *jwt.Token) (interface{}, error) {
		return typestructinterface.JwtKey, nil
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
