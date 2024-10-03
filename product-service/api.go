package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func product(r *gin.Engine, db *gorm.DB) {
	r.GET("/products", func(c *gin.Context) {
		var product []Product
		result := db.Find(&product)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		var product Product

		id := c.Param("id")
		result := db.Find(&product).Where("id = ?", id)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	})

	r.POST("/products", func(c *gin.Context) {
		var product Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		check_name := db.Find(product).Where("name = ?", product.Name)

		if check_name.RowsAffected != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Message":"Product already exist"})
			return
		}

		db.Save(&product)

		c.JSON(http.StatusOK, product)

	})

	r.PUT("/products/:id", func(c *gin.Context) {
		var product Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		var existingProduct Product

		// Find the existing product
		if err := db.First(&existingProduct, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		// Update the existing product with the new data
		if err := db.Model(&existingProduct).Updates(product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, existingProduct)

	})

	r.DELETE("/products/:id", func(c *gin.Context) {
		var product Product

		result := db.First(&product, c.Param("id"))

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		db.Delete(&product)
		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})

		db.Exec("UPDATE tasks SET id = id - 1 WHERE id > ?", c.Param("id"))
	})
}

