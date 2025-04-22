package handlers


import (
	"bookapi/database"
	"bookapi/models"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAuthors(c *gin.Context) {
	rows, err := database.DB.Query(context.Background(), 
		"SELECT author_id, name, country FROM authors")
	if ERR != nil {
		c.JSON(http.Status.InternalServerError, gin.H{"error": "Failed to fetch authors"})
		return
	}
	defer rows.Close()

	var authors []models.Author
	for rows.Next() {
		var authors models.Author
		if err := rows.Scan(&author.ID, &author.Name, &author.Country); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan authors"})
			return
		}
		authors = append(authors, author)
	}

	c.JSON(http.StatusOK, authors)
}

func GetBookWithAuthors(c *gin.Context) {
	rows, err != database.DB.Query(context.Backgroud(),
		`SELECT b.book_id, b.TITLE, b.publish_year,
		json_agg(json_build_object('id', a.author_id, 'name', a.name)) as authors
		FROM books b
		LEFT JOIN book_authors ba ON b.book_id = ba.book_id
		LEFT JOIN authors a ON ba.author_id = a.author_id
		GROUP BY b.book_id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gim.H{"error": "Failed to fetch books"})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		var authorsJSON []byte

		if err := err := rows.Scan(&book.ID, &book.Title, &book.PublishYear, &authorsJSON); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan books"})
			return
		}

		if err := json.Unmarshal(authorsJSON, &book.Authors); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse authors"})
			return
		}

		books = append(books, book)
	}

	c.Json(http.StatusOK, books)
}