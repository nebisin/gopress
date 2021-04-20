package repository

import (
	"github.com/nebisin/gopress/models"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *postRepository {
	return &postRepository{db: db}
}

// Save method takes post model and create that post
// in the database. It returns error if exist any.
func (r *postRepository) Save(p *models.Post) error {
	if err := r.db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

// FindById method find one post by given id.
func (r *postRepository) FindById(id uint) (models.Post, error) {
	var post models.Post
	if err := r.db.Preload("Author").First(&post, id).Error; err != nil {
		return models.Post{}, err
	}

	return post, nil
}

// UpdateById method update one post
// It takes old post and new post and return error if any.
func (r *postRepository) UpdateById(post *models.Post, newPost models.Post) error {
	if err := r.db.Model(&post).Updates(newPost).Error; err != nil {
		return err
	}

	return nil
}

// DeleteById method delete one post by given id.
func (r *postRepository) DeleteById(id uint) error {
	if err := r.db.Delete(&models.Post{}, id).Error; err != nil {
		return err
	}

	return nil
}

// FindMany method gets all published posts in the limits
// ordered by creation time.
// If limit is not provided it's 10 by default.
func (r *postRepository) FindMany(limit int) ([]models.Post, error) {
	if limit == 0 {
		limit = 10
	}

	var posts []models.Post
	if err := r.db.
		Limit(limit).
		Order("created_at desc").
		Preload("Author").
		Where("is_published = ?", true).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// FindPostsByUserId method gets given users posts
// just published ones.
func (r postRepository) FindPostsByUserId(uid uint) ([]models.Post, error) {
	var posts []models.Post

	if err := r.db.
		Order("created_at desc").
		Preload("Author").
		Where("author_id = ?", uid).
		Where("is_published = ?", true).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// FindMyPosts method gets given users posts
// including both published and unpublished ones.
func (r postRepository) FindMyPosts(uid uint) ([]models.Post, error) {
	var posts []models.Post

	if err := r.db.
		Order("created_at desc").
		Preload("Author").
		Where("author_id = ?", uid).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}