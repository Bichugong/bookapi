package handlers

import (
	"bookapi/config"
	"bookapi/database"
	"bookapi/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Сохраняем пользователя
	_, err = database.DB.Exec(context.Background(),
		"INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3)",
		user.Username, string(hashedPassword), user.Email)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func LoginUser(cfg config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user models.User
        if err := c.ShouldBindJSON(&user); err != nil {
            log.Println("Invalid request format:", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
            return
        }

        var storedUser models.User
        err := database.DB.QueryRow(context.Background(),
            "SELECT username, password_hash FROM users WHERE username = $1", user.Username).
            Scan(&storedUser.Username, &storedUser.PasswordHash)
        if err != nil {
            log.Println("Error fetching user:", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(user.Password)); err != nil {
            log.Println("Password mismatch:", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
            return
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
            Username: user.Username,
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            },
        })

        tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"token": tokenString})
    }
}