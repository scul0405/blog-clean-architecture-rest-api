package usecase

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/auth/mock"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAuthUseCase_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			SymmetricKey: "secret_token_symmetric_key_12345",
		},
		Logger: config.LoggerConfig{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &models.User{
		Email:    "email@gmail.com",
		Password: string(hashPassword),
	}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	defer span.Finish()

	mockAuthRepo.EXPECT().Register(ctxWithTrace, gomock.Eq(user)).Return(mockUser, nil)

	createdUserWithToken, err := authUC.Register(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUserWithToken)
	require.Nil(t, err)
}

func TestAuthUseCase_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.GetByID")
	defer span.Finish()

	mockAuthRepo.EXPECT().GetByID(ctxWithTrace, gomock.Eq(user.UserID)).Return(user, nil)

	testUser, err := authUC.GetByID(ctx, user.UserID)
	require.NoError(t, err)
	require.NotNil(t, testUser)
	require.Nil(t, err)
}

func TestAuthUseCase_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			SymmetricKey: "secret_token_symmetric_key_12345",
		},
		Logger: config.LoggerConfig{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, apiLogger)

	user := &models.LoginUser{
		Password: "123456",
		Email:    "email@gmail.com",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &models.User{
		Email:    "email@gmail.com",
		Password: string(hashPassword),
	}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.Login")
	defer span.Finish()

	mockAuthRepo.EXPECT().FindByEmail(ctxWithTrace, gomock.Eq(user.Email)).Return(mockUser, nil)

	createdUserWithToken, err := authUC.Login(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUserWithToken)
	require.Nil(t, err)
}
