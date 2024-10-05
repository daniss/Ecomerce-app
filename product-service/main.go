package main

import (

	"github.com/gin-gonic/gin"
)
type Product struct {
    ID          uint    `json:"id" gorm:"primaryKey"`
    Name        string  `json:"name" gorm:"type:varchar(255);not null"`
    Description string  `json:"description" gorm:"type:varchar(255);not null"`
    Price       float64 `json:"price" gorm:"default:0;not null"`
    Stock       int     `json:"stock" gorm:"default:0;not null"`
}

func main() {
	r := gin.Default()
	db := setupDatabase()
	product(r, db)
	r.Run("0.0.0.0:8080")
}