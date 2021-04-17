package repository

import (
	"github.com/nebisin/gopress/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Save(p *models.Post) error
	FindById(id uint) (models.Post, error)
	UpdateById(post *models.Post, newPost models.Post) error
	DeleteById(id uint) error
	FindMany(limit int) ([]models.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *postRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Save(p *models.Post) error {
	if err := r.db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepository) FindById(id uint) (models.Post, error) {
	var post models.Post
	if err := r.db.Preload("Author").First(&post, id).Error; err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (r *postRepository) UpdateById(post *models.Post, newPost models.Post) error {
	if err := r.db.Model(&post).Updates(newPost).Error; err != nil {
		return err
	}

	return nil
}

func (r *postRepository) DeleteById(id uint) error {
	if err := r.db.Delete(&models.Post{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *postRepository) FindMany(limit int) ([]models.Post, error) {
	if limit == 0 {
		limit = 10
	}

	var posts []models.Post
	if err := r.db.Limit(limit).Order("created_at desc").Preload("Author").Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}