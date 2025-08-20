package config

import (
	"database/sql"
	"fmt"
)

func CreateConnection() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=nasioC12no4 dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Connected to PostgreSQL!")
	return db
}
