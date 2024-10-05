package main

import (

	"github.com/gin-gonic/gin"
    "time"
)
type Product struct {
    ID          uint    `json:"id" gorm:"primaryKey"`
    Name        string  `json:"name" gorm:"type:varchar(255);not null"`
    Description string  `json:"description" gorm:"type:varchar(255);not null"`
    Price       float64 `json:"price" gorm:"default:0;not null"`
    Stock       int     `json:"stock" gorm:"default:0;not null"`
}

type Users struct {
    ID           uint      `gorm:"primaryKey"`
    Username     string    `json:"username" gorm:"type:varchar(40);not null"`
    PasswordHash string    `json:"password" gorm:"type:varchar(255);not null"`
    Role         string    `json:"role" gorm:"type:varchar(20)"`
    CreatedAt    time.Time `json:"created_at"`
}

func main() {
	r := gin.Default()
	db := setupDatabase()
	product(r, db)
    register(r, db)
    login(r, db)
	r.Run("0.0.0.0:8080")
}