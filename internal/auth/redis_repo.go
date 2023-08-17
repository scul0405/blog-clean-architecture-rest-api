package auth

import (
	"context"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
)

type RedisRepository interface {
	GetByIDCtx(ctx context.Context, key string) (*models.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error
	DeleteUserCtx(ctx context.Context, key string) error
}
