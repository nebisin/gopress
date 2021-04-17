package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title string `json:"title" gorm:"not null"`
	Body string `json:"body"`
	AuthorID *uint `json:"author_id" gorm:"not null"`
	Author *User `json:"author"`
}

type PostDTO struct {
	Title string `json:"title"`
	Body string `json:"body"`
}

func DTOToPost(dto PostDTO) Post {
	return Post{
		Title: dto.Title,
		Body: dto.Body,
	}
}

func PostToDTO(post Post) PostDTO {
	return PostDTO{
		Title: post.Title,
		Body: post.Body,
	}
}