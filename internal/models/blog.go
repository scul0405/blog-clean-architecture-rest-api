package models

import (
	"time"

	"github.com/google/uuid"
)

// Blog base model
type Blog struct {
	BlogID    uuid.UUID `json:"blog_id" db:"blog_id" validate:"omitempty,uuid"`
	AuthorID  uuid.UUID `json:"author_id,omitempty" db:"author_id"`
	Title     string    `json:"title" db:"title" validate:"required,gte=10"`
	Content   string    `json:"content" db:"content" validate:"required,gte=20"`
	ImageURL  *string   `json:"image_url,omitempty" db:"image_url" validate:"omitempty,lte=512,url"`
	Category  *string   `json:"category,omitempty" db:"category" validate:"omitempty,lte=10"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// BlogsList contains list of blogs
type BlogsList struct {
	TotalCount int         `json:"total_count"`
	TotalPages int         `json:"total_pages"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	HasMore    bool        `json:"has_more"`
	Blogs      []*BlogBase `json:"blogs"`
}

// BlogBase contains data when update and response to client
type BlogBase struct {
	BlogID    uuid.UUID `json:"blog_id" db:"blog_id"`
	AuthorID  uuid.UUID `json:"author_id,omitempty" db:"author_id"`
	Title     string    `json:"title" db:"title" validate:"omitempty,gte=10"`
	Content   string    `json:"content" db:"content" validate:"omitempty,gte=20"`
	ImageURL  *string   `json:"image_url,omitempty" db:"image_url" validate:"omitempty,lte=512,url"`
	Category  *string   `json:"category,omitempty" db:"category" validate:"omitempty,lte=10"`
	Author    string    `json:"author" db:"author"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
