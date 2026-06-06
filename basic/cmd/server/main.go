package main

import (
	"log"

	"basic/internal/platform/config"
	"basic/internal/platform/database"
	"basic/internal/post"
	"basic/internal/server"
	"basic/internal/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Schema is managed by the migrate command (cmd/migrate), not here. The
	// server assumes the database is already at the expexcted version.
	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	// Build each layer bottom-up and hand it to the one above. This is the
	// whole "dependency injection" idea: nothing reaches for a global, every
	// piece is given what it needs here, in one place.
	users := user.NewService(user.NewRepository(db))
	posts := post.NewService(post.NewRepository(db), users) // users satisfies post.UserChecker

	r := server.New(
		user.NewHandler(users),
		post.NewHandler(posts),
	)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
