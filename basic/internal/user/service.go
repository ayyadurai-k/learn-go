package user

import (
	"errors"
	"fmt"
	"strings"
)

var ErrEmailTaken = errors.New("email already registered")

// Service holds the business rules for users. Handlers call into it; it calls
// the repository. It never touches HTTP or the database directly.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	Name  string
	Email string
}

func (s *Service) Create(in CreateInput) (*User, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	if in.Name == "" || in.Email == "" {
		return nil, fmt.Errorf("name and email are required")
	}

	switch _, err := s.repo.GetByEmail(in.Email); {
	case err == nil:
		return nil, ErrEmailTaken
	case !errors.Is(err, ErrNotFound):
		return nil, err
	}

	u := &User{Name: in.Name, Email: in.Email}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Get(id uint) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) List() ([]User, error) {
	return s.repo.List()
}

// Exists lets other domains (posts) confirm an author without depending on the
// user repository or model internals.
func (s *Service) Exists(id uint) (bool, error) {
	switch _, err := s.repo.GetByID(id); {
	case err == nil:
		return true, nil
	case errors.Is(err, ErrNotFound):
		return false, nil
	default:
		return false, err
	}
}
