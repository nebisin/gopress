package repository

import "github.com/nebisin/gopress/models"

type PostRepository interface {
	SavePost(p *models.Post) error
	FindPostById(id uint) (models.Post, error)
	UpdatePostById(post *models.Post, newPost models.Post) error
	DeletePostById(id uint) error
	FindPosts(limit int) ([]models.Post, error)
}

func (r *Repository) SavePost(p *models.Post) error {
	if err := r.db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindPostById(id uint) (models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, id).Error; err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (r *Repository) UpdatePostById(post *models.Post, newPost models.Post) error {
	if err := r.db.Model(&post).Updates(newPost).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeletePostById(id uint) error {
	if err := r.db.Delete(&models.Post{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r Repository) FindPosts(limit int) ([]models.Post, error) {
	if limit == 0 {
		limit = 10
	}

	var posts []models.Post
	if err := r.db.Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}