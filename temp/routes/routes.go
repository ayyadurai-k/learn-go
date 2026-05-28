package routes

import (
	"product-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/products", handlers.CreateProduct)
		v1.GET("/products", handlers.GetProducts)
		v1.GET("/products/:id", handlers.GetProduct)
		v1.PUT("/products/:id", handlers.UpdateProduct)
		v1.DELETE("/products/:id", handlers.DeleteProduct)
	}

	return r
}
