package handlers

import (
	"bookapi/database"
	"bookapi/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	//Хешируем пароль
	hashedPassword, err := dcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	//Сохраняем пользователя
	_, err = database.DB.Exec(context.Background(),
		"INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3)",
		user.Username, string(hashedPassword), user.Email)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	//Ищем пользователя
	var storedUser models.User
	err := database.DB.QueryRow(context.Background(),
		"SELECT username, password_hash FROM users WHERE username = $1", user.Username).
		Scan(&storedUser.Username, &storedUser.PasswordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	//Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(user.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
	}

	//Генерируем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.now().Add(24 * time.Hour)),
		},
	}) 

	tokenString, err := token.SignedString([]byte("your_jwt_secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Faild to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}


