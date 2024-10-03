package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"fmt"
)

func product(r *gin.Engine) {
	r.GET("/products", func(c *gin.Context) {
		
	})

	r.GET("/products/:id", func(c *gin.Context) {
		
	})

	r.POST("/products", func(c *gin.Context) {

	})

	r.PUT("/products/:id", func(c *gin.Context) {

	})

	r.DELETE("/products/:id", func(c *gin.Context) {

	})
}

func 