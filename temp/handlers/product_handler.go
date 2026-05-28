package handlers

import (
	"net/http"

	"product-api/config"
	"product-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateProductInput struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Quantity    int     `json:"quantity" binding:"gte=0"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name" binding:"omitempty,min=1,max=255"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,gte=0"`
	Quantity    *int     `json:"quantity" binding:"omitempty,gte=0"`
}

func CreateProduct(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    input.Quantity,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func GetProduct(c *gin.Context) {
	var product models.Product
	if err := config.DB.First(&product, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func UpdateProduct(c *gin.Context) {
	var product models.Product
	if err := config.DB.First(&product, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product"})
		return
	}

	var input UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}
	if input.Quantity != nil {
		updates["quantity"] = *input.Quantity
	}

	if err := config.DB.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func DeleteProduct(c *gin.Context) {
	var product models.Product
	if err := config.DB.First(&product, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product"})
		return
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}
