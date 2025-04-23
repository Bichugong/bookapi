package main

import (
	"bookapi/config"
	"bookapi/database"
	"bookapi/handlers"
	"bookapi/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Загрузка конфига
	cfg := config.LoadConfig()

	// Инициализация БД
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	// Инициализация роутера
	r := gin.Default()

	// Незащищенные роуты
	r.POST("/register", handlers.RegisterUser)
	r.POST("/login", handlers.LoginUser(cfg))

	// Защищенные роуты
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware(cfg))
	{
    authGroup.GET("/authors", handlers.GetAuthors)
    authGroup.GET("/books", handlers.GetBooksWithAuthors)
	}

	// Запуск сервера
	log.Println("Server running on :8080")
	r.Run(":8080")
}