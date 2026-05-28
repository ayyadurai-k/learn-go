package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Quantity    int     `json:"quantity" binding:"gte=0"`
}
