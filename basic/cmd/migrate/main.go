// Command migrate applies or rolls back database schema migrations. It is run
// deliberately (locally or as a deploy step), separately from the API server,
// so schema changes are never tied to a server starting up.
//
//	go run ./cmd/migrate up        # apply all pending migrations
//	go run ./cmd/migrate down      # roll back the most recent one
//	go run ./cmd/migrate version   # print the current schema version
package main

import (
	"fmt"
	"log"
	"os"

	"basic/internal/platform/config"
	"basic/internal/platform/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down|version>")
	}

	switch os.Args[1] {
	case "up":
		if err := database.Up(cfg.DatabaseURL); err != nil {
			log.Fatalf("migrate up: %v", err)
		}
		fmt.Println("schema is up to date")

	case "down":
		if err := database.Down(cfg.DatabaseURL); err != nil {
			log.Fatalf("migrate down: %v", err)
		}
		fmt.Println("rolled back one migration")

	case "version":
		v, dirty, err := database.Version(cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("migrate version: %v", err)
		}
		fmt.Printf("version=%d dirty=%t\n", v, dirty)

	default:
		log.Fatalf("unknown command %q (want up, down, or version)", os.Args[1])
	}
}
