package routes

import (
	"log"
	"net/http"
)

func RegisterRouter(res http.ResponseWriter, req *http.Request) {
	log.Println("注册--原生HTTP路由")
}
