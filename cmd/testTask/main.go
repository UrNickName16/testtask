package main

import (
	"log"

	"testtask/internal/server"
)

func main() {
	srv := server.NewServer()
	log.Println("Starting server on :8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
