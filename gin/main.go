package main

import (
	"gin/middlewares"

	"github.com/gin-gonic/gin"
)


func GetData(c *gin.Context){
	c.JSON(200,gin.H{
		"data" :"data",
	})
}
func GetData1(c *gin.Context){
	c.JSON(200,gin.H{
		"data" :"data1",
	})
}
func GetData2(c *gin.Context){
	c.JSON(200,gin.H{
		"data" :"data2",
	})
}

func main()  {
	router := gin.Default()

	// router.Use(middlewares.Authenticate) // apply to all routes

	router.GET("/get-data",GetData)
	router.GET("/get-data1",GetData1)
	router.GET("/get-data2",middlewares.Authenticate,GetData2) // specific middleware

	router.Run()


}