// package config

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/lib/pq"
// )

// func CreateConnection() *sql.DB {
// 	// connStr := "postgresql://postgres:nasioC12no4@db.vvqcalobexngeftklnkc.supabase.co:5432/postgres"
// 	connStr := "postgresql://postgres:nasioC12no4321!@db.vvqcalobexngeftklnkc.supabase.co:5432/postgres"

// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}

//		fmt.Println("✅ Connected to Supabase PostgreSQL!")
//		return db
//	}
package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func CreateConnection() *sql.DB {
	dsn := os.Getenv("postgresql://postgres:nasioC12no4@db.vvqcalobexngeftklnkc.supabase.co:5432/postgres")
	if dsn == "" {
		panic("❌ DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Connected to Supabase PostgreSQL!")
	return db
}
