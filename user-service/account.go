package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	// "github.com/joho/godotenv"
)

type UsersCheck struct {
	Username string `json:"username" binding:"required"`
	PasswordHash string `json:"password" binding:"required"`
}

func HashPassword(PasswordHash string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(PasswordHash), 14)
	return string(bytes), err
}

func HashCompare(compare string, passwordhash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordhash), []byte(compare))
	return err
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("Role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "No role found"})
            c.Abort()
            return
        }

        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
        c.Abort()
    }
}

func register(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", func(c *gin.Context) {
		var user UsersCheck


		decode := json.NewDecoder(c.Request.Body)
		decode.DisallowUnknownFields()
		if err := decode.Decode(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user.Username == "" || user.PasswordHash == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Username and password are required"})
			return
		}

		var existingUser Users
		if rec := db.Where("username = ?", user.Username).First(&existingUser); rec.Error == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Username already exist"})
			return
		}

		var err error
		user.PasswordHash, err = HashPassword(user.PasswordHash)

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		var toregister Users
		toregister.Username = user.Username
		toregister.PasswordHash = user.PasswordHash
		toregister.Role = "User"
		result := db.Create(&toregister)

		if result.RowsAffected == 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "Didn't work"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User " + user.Username + " succesfuly created"})
	})
}

func createToken(user Users) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["userID"] = user.ID
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	claims["role"] = user.Role
	key := os.Getenv("SECRETKEY")
	if key == "" {
		return "", fmt.Errorf("SECRETKEY environment variable is not set")
	}
	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func login(r *gin.Engine, db *gorm.DB) {
	r.POST("/login", func(c *gin.Context) {
		var user UsersCheck

		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user.Username == "" || user.PasswordHash == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Username and password are required"})
			return
		}

		var dbUser Users
		if rec := db.Where("username = ?", user.Username).First(&dbUser); rec.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User doesn't exist"})
			return
		} else if rec.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": rec.Error.Error()})
			return
		}

		if err := HashCompare(user.PasswordHash, dbUser.PasswordHash); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong password"})
			return
		}

		token, err := createToken(dbUser)
		
		if err != nil {
			fmt.Println("Error creating token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Access Token wasn't generated"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
}
