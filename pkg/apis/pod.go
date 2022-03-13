package apis

import (
	"gin-client-go/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPods(c *gin.Context) {
	namespaceName := c.Param("namespaceName")
	pods, err := service.GetPods(namespaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, pods)
}
