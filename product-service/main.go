package main

import (
	"github.com/gin-gonic/gin"
)
type Product struct {
    ID          uint    `json:"id" gorm:"primaryKey"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
}

func main() {
	r := gin.Default()
	db := setupDatabase()
	product(r, db)

	r.Run()
}