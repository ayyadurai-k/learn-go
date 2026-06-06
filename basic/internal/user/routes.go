package user

import "github.com/gin-gonic/gin"

// Register is this domain's route table — think of it as the app-level urls.py.
// It maps URLs to handler methods; the methods themselves (the "views") live in
// handler.go. server.New calls this for each domain, like Django's
// include('users.urls').
func (h *Handler) Register(r gin.IRouter) {
	r.POST("/users", h.create)
	r.GET("/users", h.list)
	r.GET("/users/:id", h.get)
}
