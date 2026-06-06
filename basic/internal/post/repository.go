package post

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("post not found")

type Repository interface {
	Create(p *Post) error
	GetByID(id uint) (*Post, error)
	List() ([]Post, error)
	ListByUser(userID uint) ([]Post, error)
}

type gormRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

func (r *gormRepo) Create(p *Post) error {
	return r.db.Create(p).Error
}

func (r *gormRepo) GetByID(id uint) (*Post, error) {
	var p Post
	err := r.db.First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &p, err
}

func (r *gormRepo) List() ([]Post, error) {
	var posts []Post
	return posts, r.db.Order("id desc").Find(&posts).Error
}

func (r *gormRepo) ListByUser(userID uint) ([]Post, error) {
	var posts []Post
	return posts, r.db.Where("user_id = ?", userID).Order("id desc").Find(&posts).Error
}
