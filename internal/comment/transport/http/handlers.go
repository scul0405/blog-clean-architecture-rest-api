package http

import (
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/comment"
	commentAsynq "github.com/scul0405/blog-clean-architecture-rest-api/internal/comment/transport/asynq"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	asynqPkg "github.com/scul0405/blog-clean-architecture-rest-api/pkg/asynq"
	httpErrors "github.com/scul0405/blog-clean-architecture-rest-api/pkg/http_errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"net/http"
	"time"
)

type commentHandlers struct {
	cfg       *config.Config
	commentTD commentAsynq.CommentTaskDistributor
	commentUC comment.UseCase
	logger    logger.Logger
}

func NewCommentHandlers(cfg *config.Config, commentUC comment.UseCase, commentTD commentAsynq.CommentTaskDistributor, logger logger.Logger) comment.Handlers {
	return &commentHandlers{
		cfg:       cfg,
		commentUC: commentUC,
		commentTD: commentTD,
		logger:    logger,
	}
}

// Create godoc
// @Summary Create comment
// @Description create comment, returns comment
// @Tags Comment
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.Comment true "input data"
// @Success 201 {object} models.Comment
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /comments [post]
func (h *commentHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentHandlers.Create")
		defer span.Finish()

		commentReq := &models.Comment{}
		if err := utils.ReadRequest(c, commentReq); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdComment, err := h.commentUC.Create(ctx, commentReq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdComment)
	}
}

// GetByID godoc
// @Summary Get comment by id
// @Description get comment by comment_id, returns comment
// @Tags Comment
// @Accept json
// @Produce json
// @Param comment_id path string true "comment_id"
// @Success 200 {object} models.CommentBase
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /comments/{comment_id} [get]
func (h *commentHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentHandlers.GetByID")
		defer span.Finish()

		commentID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		commentByID, err := h.commentUC.GetByID(ctx, commentID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, commentByID)
	}
}

// Update godoc
// @Summary Update comment by id
// @Description update comment, returns comment
// @Tags Comment
// @Accept json
// @Produce json
// @Security Bearer
// @Param comment_id path string true "comment_id"
// @Param request body models.CommentBase true "input data"
// @Success 200 {object} models.CommentBase
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /comments/{comment_id} [patch]
func (h *commentHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentHandlers.Update")
		defer span.Finish()

		commentUID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		commentReq := &models.CommentBase{}
		if err := utils.ReadRequest(c, commentReq); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		commentReq.CommentID = commentUID

		updatedComment, err := h.commentUC.Update(ctx, commentReq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedComment)
	}
}

// Delete godoc
// @Summary Delete comment by id
// @Description Delete comment by comment_id
// @Tags Comment
// @Accept json
// @Produce json
// @Security Bearer
// @Param comment_id path string true "comment_id"
// @Success 200 {string} string "success"
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /comments/{comment_id} [delete]
func (h *commentHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentHandlers.Delete")
		defer span.Finish()

		commentUID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.commentUC.Delete(ctx, commentUID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// List godoc
// @Summary List comments by blog_id
// @Description List comments by blog_id, return list of comments
// @Tags Comment
// @Accept json
// @Produce json
// @Param blog_id query string true "blog id"
// @Success 200 {object} models.CommentsList
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /comments [get]
func (h *commentHandlers) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentHandlers.List")
		defer span.Finish()

		blogUID, err := uuid.Parse(c.QueryParam("blog_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		commentsList, err := h.commentUC.List(ctx, blogUID, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, commentsList)
	}
}

// Like godoc
// @Summary Like comment by id
// @Description like comment, returns comment
// @Tags Comment
// @Accept json
// @Produce json
// @Security Bearer
// @Param comment_id path string true "comment_id"
// @Success 200 {object} models.CommentBase
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /comments/{comment_id}/like [patch]
func (h *commentHandlers) Like() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentHandlers.Like")
		defer span.Finish()

		userUID, err := utils.GetUserUIDFromCtx(ctx)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewUnauthorizedError(err)))
		}

		commentUID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		opts := []asynq.Option{
			asynq.MaxRetry(10),
			asynq.ProcessIn(5 * time.Second),
			asynq.Queue(asynqPkg.QueueCritical),
		}

		payload := &commentAsynq.LikeCommentPayload{
			UserUID:   userUID,
			CommentID: commentUID,
		}

		err = h.commentTD.DistributeTaskLikeComment(ctx, payload, opts...)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// Dislike godoc
// @Summary Dislike comment by id
// @Description dislike comment, returns comment
// @Tags Comment
// @Accept json
// @Produce json
// @Security Bearer
// @Param comment_id path string true "comment_id"
// @Success 200 {object} models.CommentBase
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /comments/{comment_id}/dislike [patch]
func (h *commentHandlers) Dislike() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "commentHandlers.Dislike")
		defer span.Finish()

		userUID, err := utils.GetUserUIDFromCtx(ctx)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewUnauthorizedError(err)))
		}

		commentUID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		opts := []asynq.Option{
			asynq.MaxRetry(10),
			asynq.ProcessIn(5 * time.Second),
			asynq.Queue(asynqPkg.QueueCritical),
		}

		payload := &commentAsynq.DislikeCommentPayload{
			UserUID:   userUID,
			CommentID: commentUID,
		}

		err = h.commentTD.DistributeTaskDislikeComment(ctx, payload, opts...)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}
