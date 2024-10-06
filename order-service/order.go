package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"encoding/json"
	"io"
)

type CustomClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type Product struct {
    ID          uint    `json:"id" gorm:"primaryKey"`
    Name        string  `json:"name" gorm:"type:varchar(255);not null"`
    Description string  `json:"description" gorm:"type:varchar(255);not null"`
    Price       float64 `json:"price" gorm:"default:0;not null"`
    Stock       int     `json:"stock" gorm:"default:0;not null"`
}

func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No or invalid authorization given"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		key := os.Getenv("SECRETKEY")
		if key == "" {
			c.Abort()
			panic("SECRETKEY environment variable is not set")
		}
		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("Role", claims.Role)
		c.Set("UserID", claims.UserID)
		c.Set("Authorization", tokenString)

		c.Next()
	}
}

func order(r *gin.Engine, db *gorm.DB) {
	r.GET("/orders", jwtAuthMiddleware(), func(c *gin.Context) {
		var orders []Order
		result := db.Where("user_id = ?", c.MustGet("UserID")).Find(&orders)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No orders found on this account"})
			return
		}

		c.JSON(http.StatusOK, orders)
	})

	r.GET("/orders/:id", jwtAuthMiddleware(), func(c *gin.Context) {
		var order Order
		id := c.Param("id")

		result := db.Where("id = ?", id, "user_id = ?", c.MustGet("UserID")).First(&order)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found or you don't have permission to view this order"})
			return
		}

		c.JSON(http.StatusOK, order)
	})

	r.POST("/orders", jwtAuthMiddleware(), func(c *gin.Context) {
		var order Order

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"service": "service",
		})

		tokenString, _ := token.SignedString([]byte(os.Getenv("SECRETKEY")))

		var order_id string = strconv.FormatUint(uint64(order.ID), 10)

		req, err := http.NewRequest(http.MethodGet, "http://product-service:8080/products/" + order_id, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}
		req.Header.Set("Authorization", "Bearer " + tokenString)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to notify external service"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Product does not exist"})
			return
		}

		body, err := io.ReadAll(resp.Body)
    	if err != nil {
    	    fmt.Println("Error reading response body:", err)
    	    return
    	}

    	var product Product
    	err = json.Unmarshal(body, &product)
    	if err != nil {
    	    fmt.Println("Error unmarshaling JSON:", err)
    	    return
    	}

		if product.Stock < order.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock"})
			return
		}
		
		if err := db.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}


		// data := map[string]interface{}{
		// 	"order_id": order.ID,
		// 	"product_id": order.ProductID,
		// 	"quantity": order.Quantity,
		// }
		// resp, err := http.NewRequest(http.MethodPut, "http://product-service:8080/products/", "application/json")
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to notify external service"})
		// 	return
		// }
		// defer resp.Body.Close()

		// if resp.StatusCode != http.StatusOK {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "External service returned non-OK status"})
		// 	return
		// }

		c.JSON(http.StatusOK, order)
	})

	r.PUT("/orders/:id", jwtAuthMiddleware(), func(c *gin.Context) {
		var order Order
		id := c.Param("id")

		result := db.Where("id = ?", id, "user_id = ?", c.MustGet("UserID")).First(&order)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found or you don't have permission to update this order"})
			return
		}

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	})

	r.DELETE("/orders/:id", jwtAuthMiddleware(), func(c *gin.Context) {
		var order Order
		id := c.Param("id")

		result := db.Where("id = ?", id, "user_id = ?", c.MustGet("UserID")).First(&order)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found or you don't have permission to delete this order"})
			return
		}

		if err := db.Delete(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order " + id + " deleted"})
	})
}