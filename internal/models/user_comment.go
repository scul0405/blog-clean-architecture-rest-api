package models

import (
	"github.com/google/uuid"
	"time"
)

// UserComments model
type UserComments struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id" validate:"omitempty,uuid"`
	CommentID uuid.UUID `json:"comment_id" db:"comment_id" validate:"omitempty,uuid"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
