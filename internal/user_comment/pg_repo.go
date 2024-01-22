//go:generate mockgen -source pg_repo.go -destination mock/pg_repo_mock.go -package mock
package user_comment

import (
	"context"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
)

type Repository interface {
	GetByID(ctx context.Context, userComment *models.UserComments) error
	Create(ctx context.Context, userComment *models.UserComments) error
	Delete(ctx context.Context, userComment *models.UserComments) error
}
