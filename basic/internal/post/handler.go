package post

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) create(c *gin.Context) {
	var req struct {
		UserID uint   `json:"user_id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	p, err := h.svc.Create(CreateInput{UserID: req.UserID, Title: req.Title, Body: req.Body})
	switch {
	case errors.Is(err, ErrUnknownUser):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	case err != nil:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusCreated, p)
	}
}

func (h *Handler) get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	p, err := h.svc.Get(uint(id))
	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	default:
		c.JSON(http.StatusOK, p)
	}
}

func (h *Handler) list(c *gin.Context) {
	// Optional ?user_id=N filter, handled in the service.
	var userID uint
	if raw := c.Query("user_id"); raw != "" {
		n, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		userID = uint(n)
	}

	posts, err := h.svc.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	c.JSON(http.StatusOK, posts)
}
