package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Order struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	r := gin.Default()
	db := setupDatabase()
	order(r, db)
    
	r.Run("0.0.0.0:8080")
}