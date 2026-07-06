package database

import (
	"database/sql"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(DBPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return nil, err
	}

	if err := readSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func readSchema(db *sql.DB) error {
	bytes, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(bytes))
	if err != nil {
		return err
	}
	return nil
}
