package main

import (
	"fmt"
	"gin-client-go/pkg/config"
	"gin-client-go/pkg/router"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

func main() {
	engine := gin.Default()
	gin.SetMode(gin.DebugMode)
	router.InitRouter(engine)
	err := engine.Run(fmt.Sprintf("%s:%d", config.GetString(config.ServerHost), config.GetInt(config.ServerPort)))
	if err != nil {
		klog.Fatal(err)
		return
	}

}
