package main

import (
	"log"
	"net/http"

	"example.com/filmserver/internal/server"
)

func main() {
	srv := server.New()
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
