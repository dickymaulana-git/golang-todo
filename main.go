package main

import (
	"fmt"
	"log"
	"main/routers"
	"net/http"

	_ "github.com/lib/pq" // Importing pq for PostgreSQL driver
)

func main() {
	r := routers.Router()
	fmt.Println("Server is running on port 8000")

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: r,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
