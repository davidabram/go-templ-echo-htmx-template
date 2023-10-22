package db

import (
	"database/sql"
	"os"
	"fmt"

	_ "github.com/libsql/libsql-client-go/libsql"
)

type Database struct {
	db *sql.DB
}

func (d *Database) New() (error) {
	db, err := sql.Open("libsql", os.Getenv("DATABASE_URL"))

	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("Error pinging database: %v", err)
	}

	d.db = db
	return nil
}
