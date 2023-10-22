package db

import (
	"database/sql"
)

type Author struct {
	AuthorID int
	Name     string
	Bio      sql.NullString // Using NullString to handle SQL NULL values
}

type Category struct {
	CategoryID int
	Name       string
}

type Book struct {
	BookID          int
	Title           string
	PublicationYear int
	AuthorID        int
	CategoryID      int
	Price           float64
}

type FullBookInfo struct {
	Book
	Author   Author
	Category Category
}
