package middlewares

import "github.com/gin-gonic/gin"

func Authenticate(c *gin.Context) {

	token := c.Request.Header.Get("Token")

	if token == "" {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "Token not present",
		})
		return
	}

	if token != "auth" {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "Invalid Token",
		})
		return
	}

	c.Next()
}
