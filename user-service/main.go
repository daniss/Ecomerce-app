package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

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
    register(r, db)
    login(r, db)
	r.Run("0.0.0.0:8080")
}