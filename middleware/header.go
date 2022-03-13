package middleware

import "github.com/gin-gonic/gin"

func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "origin, content-type, authorization, accept")
		c.Header("Allow", "GET, POST, PUT, DELETE, OPTIONS, PATCH, Header")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(200)
	}
}
