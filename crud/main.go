package main

import (
	"crud/controllers"
	"crud/initializers"

	"github.com/gin-gonic/gin"
)
    

func init(){
    initializers.LoadEnvVariables()
    initializers.ConnectDB()
}

func main() {
  router := gin.Default()
  router.POST("/posts",controllers.PostCreate)
  router.GET("/posts/:id",controllers.PostDetail)

  router.Run() // listens on 8080 by default
}