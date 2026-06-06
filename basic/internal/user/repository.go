package user

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("user not found")

// Repository is the storage boundary for users. The service depends on this
// interface rather than on *gorm.DB, which is what lets us hand it a fake in
// tests (see post/service_test.go for the same idea in action).
type Repository interface {
	Create(u *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	List() ([]User, error)
}

type gormRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

func (r *gormRepo) Create(u *User) error {
	return r.db.Create(u).Error
}

func (r *gormRepo) GetByID(id uint) (*User, error) {
	var u User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *gormRepo) GetByEmail(email string) (*User, error) {
	var u User
	err := r.db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *gormRepo) List() ([]User, error) {
	var users []User
	return users, r.db.Find(&users).Error
}
