package http

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/auth"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	httpErrors "github.com/scul0405/blog-clean-architecture-rest-api/pkg/http_errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"io"
	"net/http"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	logger logger.Logger
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, logger logger.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, logger: logger}
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.User true "input data"
// @Success 201 {object} models.User
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/register [post]
func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "authHandlers.Register")
		defer span.Finish()

		user := &models.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userWithToken, err := h.authUC.Register(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, userWithToken)
	}
}

// GetByID godoc
// @Summary Get user
// @Description Get user by user's id, return user
// @Tags Auth
// @Accept json
// @Param id path string true "id"
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/{id} [get]
func (h *authHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "authHandlers.GetByID")
		defer span.Finish()

		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user, err := h.authUC.GetByID(ctx, userID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Login godoc
// @Summary Login user
// @Description login user, returns user and access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.User true "input data"
// @Success 200 {object} models.User
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/login [post]
func (h *authHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "authHandlers.Login")
		defer span.Finish()

		user := &models.LoginUser{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userWithToken, err := h.authUC.Login(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, userWithToken)
	}
}

// UploadAvatar godoc
// @Summary Upload avatar user
// @Description upload avatar user, returns user
// @Tags Auth
// @Accept json
// @Produce json
// @Param file formData file  true "avatar"
// @Param id path string true "user id"
// @Param bucket query string true "minio bucket"
// @Success 200 {object} models.User
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/{id}/avatar [post]
func (h *authHandlers) UploadAvatar() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "authHandlers.UploadAvatar")
		defer span.Finish()

		bucket := c.QueryParam("bucket")
		uID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		image, err := utils.ReadImage(c, "file")
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		file, err := image.Open()
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		defer file.Close()

		binaryImage := bytes.NewBuffer(nil)
		if _, err = io.Copy(binaryImage, file); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		contentType, err := utils.CheckImageFileContentType(binaryImage.Bytes())
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		reader := bytes.NewReader(binaryImage.Bytes())

		updatedUser, err := h.authUC.UploadAvatar(ctx, uID, models.UploadInput{
			File:        reader,
			Name:        image.Filename,
			Size:        image.Size,
			ContentType: contentType,
			BucketName:  bucket,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}
