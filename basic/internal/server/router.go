package server

import (
	"net/http"

	"basic/internal/post"
	"basic/internal/user"

	"github.com/gin-gonic/gin"
)

// New builds the HTTP engine and mounts every domain's routes under /api/v1.
// Adding a new domain is a one-liner here once its handler exposes Register.
func New(users *user.Handler, posts *post.Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	users.Register(api)
	posts.Register(api)

	return r
}
