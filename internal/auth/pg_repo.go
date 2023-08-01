package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
)

type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
