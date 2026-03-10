package main

import (
	"log"
	"net/http"

	"example.com/filmserver/internal/server"
	"example.com/filmserver/pkg/env"
)

func main() {
	port := env.GetEnv("port", "8080")
	srv := server.New()
	log.Printf("Starting server on :%v", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
