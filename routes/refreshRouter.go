package routes

import (
	"log"
	"net/http"
	typestructinterface "server-by-chat/typeStructInterface"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func RefreshRouter(res http.ResponseWriter, req *http.Request) {
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
	jwtTokenResponseClaimsStruct := &typestructinterface.JwtTokenResponseClaimsStruct{}
	okToken, err := jwt.ParseWithClaims(
		httpCookie.Value,
		jwtTokenResponseClaimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			return typestructinterface.JwtKey, nil
		},
	)
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

	if time.Until(jwtTokenResponseClaimsStruct.RegisteredClaims.ExpiresAt.Local()) > 30*time.Second {
		res.WriteHeader(http.StatusBadRequest)
		log.Println("jwtToken到期还有很长时间\n\r")
		return
	}

	log.Println("开始颁发新的jwttoken\n\r")
	expirationTime := time.Now().Add(5 * time.Minute)
	log.Println("逾期时间对象\n\r", expirationTime)
	jwtTokenResponseClaimsStruct.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtTokenResponseClaimsStruct)
	tokenString, err := token.SignedString(typestructinterface.JwtKey)
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
