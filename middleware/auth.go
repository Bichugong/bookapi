package middleware

import (
    "bookapi/config" 
    "bookapi/models"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "fmt" 
)

func AuthMiddleware(cfg config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("AuthMiddleware triggered") // Отладочная печать
        authHeader := c.GetHeader("Authorization")
        fmt.Printf("Authorization Header: %s\n", authHeader) // Отладочная печать
        if authHeader == "" {
            fmt.Println("Authorization header missing") // Отладочная печать
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
        fmt.Printf("Token String: %s\n", tokenString) // Отладочная печать
        if tokenString == authHeader {
            fmt.Println("Bearer token missing") // Отладочная печать
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }

        claims := &models.Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            fmt.Printf("Using JWT Secret: %s\n", cfg.JWTSecret) // Отладочная печать
            return []byte(cfg.JWTSecret), nil
        })

        if err != nil {
            fmt.Printf("Token Parsing Error: %v\n", err) // Отладочная печать
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        if !token.Valid {
            fmt.Println("Token Invalid") // Отладочная печать
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        fmt.Printf("Token Claims: %+v\n", claims) // Отладочная печать
        c.Set("username", claims.Username)
        c.Next()
        fmt.Println("AuthMiddleware passed") // Отладочная печать
    }
}