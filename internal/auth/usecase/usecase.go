package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/auth"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
)

type authUseCase struct {
	cfg      *config.Config
	authRepo auth.Repository
	logger   logger.Logger
}

func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, logger logger.Logger) auth.UseCase {
	return &authUseCase{cfg: cfg, authRepo: authRepo, logger: logger}
}

func (u *authUseCase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	defer span.Finish()

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *authUseCase) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.GetByID")
	defer span.Finish()

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
