package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
)

type Database struct {
	Instance *sql.DB
}

func (d *Database) New() error {
	db, err := sql.Open("libsql", os.Getenv("DATABASE_URL"))

	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("Error pinging database: %v", err)
	}

	d.Instance = db
	return nil
}

func (d *Database) ListFullBookInfo() ([]FullBookInfo, error) {
	rows, err := d.Instance.Query(`
		SELECT b.*, a.author_id, a.name, a.bio, c.category_id, c.name
		FROM books b
		JOIN authors a ON b.author_id = a.author_id
		JOIN categories c ON b.category_id = c.category_id
	`)

	if err != nil {
		return nil, fmt.Errorf("Error querying database: %v", err)
	}

	defer rows.Close()

	var books []FullBookInfo
	for rows.Next() {
		var b Book
		var a Author
		var c Category
		if err := rows.Scan(
			&b.BookID, &b.Title, &b.PublicationYear,
			&b.AuthorID, &b.CategoryID, &b.Price,
			&a.AuthorID, &a.Name, &a.Bio,
			&c.CategoryID, &c.Name,
		); err != nil {
			return nil, err
		}
		books = append(books, FullBookInfo{Book: b, Author: a, Category: c})
	}

	return books, nil
}
