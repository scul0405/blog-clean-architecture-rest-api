package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/comment"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	httpErrors "github.com/scul0405/blog-clean-architecture-rest-api/pkg/http_errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"net/http"
)

type commentUseCase struct {
	cfg         *config.Config
	commentRepo comment.Repository
	logger      logger.Logger
}

func NewCommentUseCase(cfg *config.Config, commentRepo comment.Repository, logger logger.Logger) comment.UseCase {
	return &commentUseCase{cfg: cfg, commentRepo: commentRepo, logger: logger}
}

func (u *commentUseCase) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUC.Create")
	defer span.Finish()

	userUID, err := utils.GetUserUIDFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "commentUC.Create.GetUserUIDFromCtx"))
	}

	comment.AuthorID = userUID
	createdBlog, err := u.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return createdBlog, nil
}

func (u *commentUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.CommentBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUC.GetByID")
	defer span.Finish()

	comment, err := u.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (u *commentUseCase) Update(ctx context.Context, comment *models.CommentBase) (*models.CommentBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUC.Update")
	defer span.Finish()

	commentByID, err := u.commentRepo.GetByID(ctx, comment.CommentID)
	if err != nil {
		return nil, err
	}

	if err := utils.ValidateIsOwner(ctx, commentByID.AuthorID.String(), u.logger); err != nil {
		return nil, httpErrors.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "commentUC.Update.ValidateIsOwner"))
	}

	updatedBlog, err := u.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return updatedBlog, nil
}

func (u *commentUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUC.Delete")
	defer span.Finish()

	commentByID, err := u.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := utils.ValidateIsOwner(ctx, commentByID.AuthorID.String(), u.logger); err != nil {
		return httpErrors.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "commentUC.Delete.ValidateIsOwner"))
	}

	return u.commentRepo.Delete(ctx, id)
}

func (u *commentUseCase) List(ctx context.Context, blogID uuid.UUID, pq *utils.PaginationQuery) (*models.CommentsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUC.List")
	defer span.Finish()

	return u.commentRepo.List(ctx, blogID, pq)
}