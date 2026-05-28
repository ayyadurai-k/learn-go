package main

import (
	"log"

	"product-api/config"
	"product-api/routes"
)

func main() {
	config.ConnectDatabase()

	r := routes.SetupRouter()

	log.Println("server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
