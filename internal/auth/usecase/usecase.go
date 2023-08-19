package usecase

import (
	"context"
	"fmt"
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

const (
	basePrefix    = "api-auth"
	cacheDuration = 3600
)

type authUseCase struct {
	cfg       *config.Config
	authRepo  auth.Repository
	redisRepo auth.RedisRepository
	minioRepo auth.MinioRepository
	logger    logger.Logger
}

func NewAuthUseCase(
	cfg *config.Config,
	authRepo auth.Repository,
	redisRepo auth.RedisRepository,
	minioRepo auth.MinioRepository,
	logger logger.Logger) auth.UseCase {
	return &authUseCase{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo, minioRepo: minioRepo, logger: logger}
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

	cachedUser, err := u.redisRepo.GetByIDCtx(ctx, u.generateUserKey(userID.String()))
	if err != nil {
		u.logger.Errorf("authUC.GetByID: GetByIDCtx: %v", err)
	}

	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetUserCtx(ctx, u.generateUserKey(userID.String()), cacheDuration, user); err != nil {
		u.logger.Errorf("authUC.GetByID.SetUserCtx: %v", err)
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

func (u *authUseCase) UploadAvatar(ctx context.Context, userID uuid.UUID, file models.UploadInput) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.UploadAvatar")
	defer span.Finish()

	uploadInfo, err := u.minioRepo.PutObject(ctx, file)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.UploadAvatar.PutObject"))
	}

	avatarURL := u.generateMinioURL(file.BucketName, uploadInfo.Key)

	updatedUser, err := u.authRepo.Update(ctx, &models.User{
		UserID: userID,
		Avatar: &avatarURL,
	})
	if err != nil {
		return nil, err
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

func (u *authUseCase) generateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}

func (u *authUseCase) generateMinioURL(bucket string, key string) string {
	return fmt.Sprintf("%s/minio/%s/%s", u.cfg.Minio.MinioEndpoint, bucket, key)
}
