package main

import (
	"crud/initializers"
	"github.com/gin-gonic/gin"
)
    

func init(){
    initializers.LoadEnvVariables()
    initializers.ConnectDB()
}

func main() {
  router := gin.Default()
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })
  router.Run() // listens on 8080 by default
}