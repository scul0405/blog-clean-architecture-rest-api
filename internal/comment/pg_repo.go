//go:generate mockgen -source pg_repo.go -destination mock/pg_repo_mock.go -package mock
package comment

import (
	"context"
	"github.com/google/uuid"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
)

type Repository interface {
	Create(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.CommentBase, error)
	Update(ctx context.Context, comment *models.CommentBase) (*models.CommentBase, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, blogID uuid.UUID, pq *utils.PaginationQuery) (*models.CommentsList, error)
}
