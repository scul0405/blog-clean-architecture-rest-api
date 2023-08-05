package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/auth"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	httpErrors "github.com/scul0405/blog-clean-architecture-rest-api/pkg/http_errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/paseto"
)

type authUseCase struct {
	cfg      *config.Config
	authRepo auth.Repository
	logger   logger.Logger
}

func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, logger logger.Logger) auth.UseCase {
	return &authUseCase{cfg: cfg, authRepo: authRepo, logger: logger}
}

func (u *authUseCase) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	defer span.Finish()

	if err := user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := paseto.GeneratePASETOToken(createdUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.GeneratePASETOToken"))
	}

	createdUser.SanitizePassword()
	return &models.UserWithToken{
		User:        createdUser,
		AccessToken: token,
	}, nil
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

func (u *authUseCase) Login(ctx context.Context, loginReq *models.LoginUser) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Login")
	defer span.Finish()

	user, err := u.authRepo.FindByEmail(ctx, loginReq.Email)
	if err != nil {
		return nil, err
	}

	if err = user.ComparePassword(loginReq.Password); err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.Login.ComparePasswords"))
	}

	token, err := paseto.GeneratePASETOToken(user, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Login.GeneratePASETOToken"))
	}

	user.SanitizePassword()
	return &models.UserWithToken{
		User:        user,
		AccessToken: token,
	}, nil
}
