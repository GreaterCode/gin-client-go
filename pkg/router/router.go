package router

import (
	"gin-client-go/middleware"
	"gin-client-go/pkg/apis"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	middleware.InitMiddleware(r)
	r.GET("/ping", apis.Ping)
	r.GET("/namespaces", apis.GetNamespaces)
	r.GET("/namespace/:namespaceName/pods", apis.GetPods)
	//r.GET("/namespaces/:namespaceName/pod/:podName/container/:containerName", apis.ExecContainer)
}
