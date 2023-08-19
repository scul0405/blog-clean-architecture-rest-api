//go:generate mockgen -source pg_repo.go -destination mock/pg_repo_mock.go -package mock
package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
)

type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
}
