package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/blog"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	httpErrors "github.com/scul0405/blog-clean-architecture-rest-api/pkg/http_errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"net/http"
)

type blogHandlers struct {
	cfg    *config.Config
	blogUC blog.UseCase
	logger logger.Logger
}

func NewBlogHandlers(cfg *config.Config, blogUC blog.UseCase, logger logger.Logger) blog.Handlers {
	return &blogHandlers{
		cfg:    cfg,
		blogUC: blogUC,
		logger: logger,
	}
}

// Create godoc
// @Summary Create blog
// @Description create blog, returns blog
// @Tags Blog
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.Blog true "input data"
// @Success 201 {object} models.BlogBase
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /blogs [post]
func (h *blogHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "blogHandlers.Create")
		defer span.Finish()

		blogReq := &models.Blog{}
		if err := utils.ReadRequest(c, blogReq); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdBlog, err := h.blogUC.Create(ctx, blogReq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdBlog)
	}
}

// GetByID godoc
// @Summary Get blog by id
// @Description get blog by blog_id, returns blog
// @Tags Blog
// @Accept json
// @Produce json
// @Param blog_id path string true "blog_id"
// @Success 200 {object} models.BlogBase
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /blogs/{blog_id} [get]
func (h *blogHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "blogHandlers.GetByID")
		defer span.Finish()

		blogID, err := uuid.Parse(c.Param("blog_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		blogByID, err := h.blogUC.GetByID(ctx, blogID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, blogByID)
	}
}

// Update godoc
// @Summary Update blog by id
// @Description update blog, returns blog
// @Tags Blog
// @Accept json
// @Produce json
// @Security Bearer
// @Param blog_id path string true "blog_id"
// @Param request body models.BlogBase true "input data"
// @Success 200 {object} models.BlogBase
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /blogs/{blog_id} [patch]
func (h *blogHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "blogHandlers.Update")
		defer span.Finish()

		blogID, err := uuid.Parse(c.Param("blog_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		blogReq := &models.BlogBase{}
		if err := utils.ReadRequest(c, blogReq); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		blogReq.BlogID = blogID

		updatedBlog, err := h.blogUC.Update(ctx, blogReq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedBlog)
	}
}

// Delete godoc
// @Summary Delete blog by id
// @Description Delete blog by blog_id
// @Tags Blog
// @Accept json
// @Produce json
// @Security Bearer
// @Param blog_id path string true "blog_id"
// @Success 200 {string} string "success"
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /blogs/{blog_id} [delete]
func (h *blogHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "blogHandlers.Delete")
		defer span.Finish()

		blogID, err := uuid.Parse(c.Param("blog_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.blogUC.Delete(ctx, blogID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// List godoc
// @Summary List blogs
// @Description List blogs, return list of blogs
// @Tags Blog
// @Accept json
// @Produce json
// @Success 200 {object} models.BlogsList
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /blogs [get]
func (h *blogHandlers) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "blogHandlers.List")
		defer span.Finish()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		blogsList, err := h.blogUC.List(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, blogsList)
	}
}
