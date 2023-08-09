package blog

import (
	"context"
	"github.com/google/uuid"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
)

type UseCase interface {
	Create(ctx context.Context, blog *models.Blog) (*models.BlogBase, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.BlogBase, error)
	Update(ctx context.Context, blog *models.BlogBase) (*models.BlogBase, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, pq *utils.PaginationQuery) (*models.BlogsList, error)
}
