package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func product(r *gin.Engine, db *gorm.DB) {
    r.GET("/products", jwtAuthMiddleware(), func(c *gin.Context) {
        var products []Product
        result := db.Find(&products)

        if result.Error != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
            return
        }

        c.JSON(http.StatusOK, products)
    })

    r.GET("/products/:id", jwtAuthMiddleware(), func(c *gin.Context) {
        var product Product
        id := c.Param("id")

        result := db.Where("id = ?", id).First(&product)

        if result.Error != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
            return
        }

        c.JSON(http.StatusOK, product)
    })

    r.POST("/products", jwtAuthMiddleware(), func(c *gin.Context) {
        var product Product

        if err := c.BindJSON(&product); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var existingProduct Product
        checkName := db.Where("name = ?", product.Name).First(&existingProduct)

        if checkName.RowsAffected > 0 {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Product already exists"})
            return
        }

        if err := db.Save(&product).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, product)
    })

    r.PUT("/products/:id", jwtAuthMiddleware(), func(c *gin.Context) {
        var product Product

        if err := c.BindJSON(&product); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        id := c.Param("id")
        var existingProduct Product

        if err := db.First(&existingProduct, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
            return
        }

        if err := db.Model(&existingProduct).Updates(product).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, existingProduct)
    })

    r.DELETE("/products/:id", jwtAuthMiddleware(), func(c *gin.Context) {
        var product Product
        id := c.Param("id")

        if err := db.Exec("UPDATE tasks SET id = id - 1 WHERE id > ?", id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        result := db.First(&product, id)

        if result.Error != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
            return
        }

        if err := db.Delete(&product).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
    })
}