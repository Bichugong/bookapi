package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	// Настройка пула соединений
	config, err := pgxpool.ParseConfig("postgres://postgres:Zexecmdirjkt8@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	config.MaxConns = 20  // Лимит соединений

	// Подключаемся
	db, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(fmt.Sprintf("Ошибка подключения: %v", err))
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/authors", getAuthors)
	r.GET("/books", getBooksWithAuthors)

	fmt.Println("Сервер запущен на :8080")
	r.Run(":8080")
}

func getAuthors(c *gin.Context) {
	rows, err := db.Query(context.Background(), "SELECT author_id, name, country FROM authors")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var authors []map[string]interface{}
	for rows.Next() {
		var id int
		var name, country string
		if err := rows.Scan(&id, &name, &country); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		authors = append(authors, gin.H{
			"id":      id,
			"name":    name,
			"country": country,
		})
	}
	c.JSON(200, authors)
}

func getBooksWithAuthors(c *gin.Context) {
	rows, err := db.Query(context.Background(), `
		SELECT b.book_id, b.title, b.publish_year, 
		       json_agg(json_build_object('id', a.author_id, 'name', a.name)) as authors
		FROM books b
		LEFT JOIN book_authors ba ON b.book_id = ba.book_id
		LEFT JOIN authors a ON ba.author_id = a.author_id
		GROUP BY b.book_id
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []map[string]interface{}
	for rows.Next() {
		var id int
		var title string
		var year *int
		var authors []byte // Для json_agg

		if err := rows.Scan(&id, &title, &year, &authors); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		books = append(books, gin.H{
			"id":          id,
			"title":       title,
			"publish_year": year,
			"authors":     string(authors),
		})
	}
	c.JSON(200, books)
}