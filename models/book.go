package models

type Author struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type Book struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	PublishYear *int    `json:"publish_year"`
	ISBN        string  `json:"isbn"`
	Authors     []Author `json:"authors"`
}