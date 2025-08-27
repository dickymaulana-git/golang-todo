package main

import (
	"fmt"
	"log"
	"main/routers"
	"net/http"
	"os"
)

func main() {
	r := routers.Router()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("🚀 Server running on port", port)

	// server := http.Server{
	// 	Addr:    ":" + port, // listen di semua interface
	// 	Handler: r,
	// }
	server := http.Server{
		Addr:    "localhost:8000",
		Handler: r,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
