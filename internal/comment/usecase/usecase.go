package usecase

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/comment"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/user_comment"
	httpErrors "github.com/scul0405/blog-clean-architecture-rest-api/pkg/http_errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"net/http"
)

type commentUseCase struct {
	cfg             *config.Config
	commentRepo     comment.Repository
	userCommentRepo user_comment.Repository
	logger          logger.Logger
}

func NewCommentUseCase(
	cfg *config.Config,
	commentRepo comment.Repository,
	userCommentRepo user_comment.Repository,
	logger logger.Logger) comment.UseCase {
	return &commentUseCase{cfg: cfg, commentRepo: commentRepo, userCommentRepo: userCommentRepo, logger: logger}
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

func (u *commentUseCase) Like(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUC.Like")
	defer span.Finish()

	// check if user has liked comment
	userUID, err := utils.GetUserUIDFromCtx(ctx)
	if err != nil {
		return httpErrors.NewUnauthorizedError(errors.WithMessage(err, "commentUC.Like.GetUserUIDFromCtx"))
	}

	uc := &models.UserComments{UserID: userUID, CommentID: id}
	err = u.userCommentRepo.GetByID(ctx, uc)
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = u.userCommentRepo.Create(ctx, uc)
	if err != nil {
		return err
	}

	return nil
}

func (u *commentUseCase) Dislike(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUC.Dislike")
	defer span.Finish()

	// check if user has liked comment
	userUID, err := utils.GetUserUIDFromCtx(ctx)
	if err != nil {
		return httpErrors.NewUnauthorizedError(errors.WithMessage(err, "commentUC.Dislike.GetUserUIDFromCtx"))
	}

	uc := &models.UserComments{UserID: userUID, CommentID: id}
	err = u.userCommentRepo.GetByID(ctx, uc)
	if errors.Is(err, sql.ErrNoRows) || err != nil {
		return nil
	}

	err = u.userCommentRepo.Delete(ctx, uc)
	if err != nil {
		return err
	}

	return nil
}
