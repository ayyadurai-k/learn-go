package post

import (
	"errors"
	"fmt"
	"strings"
)

var ErrUnknownUser = errors.New("author does not exist")

// UserChecker is the only thing posts need from the user domain. Declaring the
// interface here (on the consumer side) instead of importing user.Service keeps
// the two domains decoupled and avoids an import cycle. user.Service happens to
// satisfy it, so main.go can pass it straight in.
type UserChecker interface {
	Exists(id uint) (bool, error)
}

type Service struct {
	repo  Repository
	users UserChecker
}

func NewService(repo Repository, users UserChecker) *Service {
	return &Service{repo: repo, users: users}
}

type CreateInput struct {
	UserID uint
	Title  string
	Body   string
}

func (s *Service) Create(in CreateInput) (*Post, error) {
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	ok, err := s.users.Exists(in.UserID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrUnknownUser
	}

	p := &Post{UserID: in.UserID, Title: in.Title, Body: in.Body}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Get(id uint) (*Post, error) {
	return s.repo.GetByID(id)
}

func (s *Service) List(userID uint) ([]Post, error) {
	if userID != 0 {
		return s.repo.ListByUser(userID)
	}
	return s.repo.List()
}
