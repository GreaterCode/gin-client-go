package router

import (
	"gin-client-go/pkg/apis"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.GET("/ping", apis.Ping)
	r.GET("/namespaces", apis.GetNamespaces)
}
