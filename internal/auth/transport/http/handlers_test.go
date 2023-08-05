package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/auth/mock"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/converter"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandlers_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock.NewMockUseCase(ctrl)

	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	authHandlers := NewAuthHandlers(cfg, mockAuthUC, apiLogger)

	gender := "male"
	user := &models.User{
		FirstName: "Liem",
		LastName:  "Le",
		Email:     "liemledeptrai@gmail.com",
		Password:  "123456",
		Gender:    &gender,
	}

	buf, err := converter.ToBytesBuffer(user)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, buf)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/api/v1/auth/register", strings.NewReader(buf.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "auth.Register")
	defer span.Finish()

	userWithToken := &models.UserWithToken{
		User: user,
	}

	mockAuthUC.EXPECT().Register(ctxWithTrace, gomock.Eq(user)).Return(userWithToken, nil)

	handlerFunc := authHandlers.Register()
	err = handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestAuthHandlers_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock.NewMockUseCase(ctrl)

	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	authHandlers := NewAuthHandlers(cfg, mockAuthUC, apiLogger)

	gender := "male"
	userUID := uuid.New()
	user := &models.User{
		UserID:    userUID,
		FirstName: "Liem",
		LastName:  "Le",
		Email:     "liemledeptrai@gmail.com",
		Password:  "123456",
		Gender:    &gender,
	}

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/auth/:id", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userUID.String())
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "auth.GetByID")
	defer span.Finish()

	mockAuthUC.EXPECT().GetByID(ctxWithTrace, gomock.Eq(user.UserID)).Return(user, nil)

	handlerFunc := authHandlers.GetByID()
	err := handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}
