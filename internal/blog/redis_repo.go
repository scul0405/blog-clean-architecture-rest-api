package blog

import (
	"context"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
)

type RedisRepository interface {
	GetBlogByIDCtx(ctx context.Context, key string) (*models.BlogBase, error)
	SetBlogCtx(ctx context.Context, key string, seconds int, blog *models.BlogBase) error
	DeleteBlogCtx(ctx context.Context, key string) error
}
