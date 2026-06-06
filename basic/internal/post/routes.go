package post

import "github.com/gin-gonic/gin"

// Register is this domain's route table (its urls.py). The handler methods it
// points at live in handler.go (the views).
func (h *Handler) Register(r gin.IRouter) {
	r.POST("/posts", h.create)
	r.GET("/posts", h.list)
	r.GET("/posts/:id", h.get)
}
