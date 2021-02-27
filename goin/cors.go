package goin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", c.GetHeader("origin")) // 允许跨域的域名
	//c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization, Token")
	//c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	//c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	//c.Header("Access-Control-Allow-Credentials", "true") // 是否允许发送Cookie

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.Next()
}
