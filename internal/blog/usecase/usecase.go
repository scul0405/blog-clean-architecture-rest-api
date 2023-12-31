package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/blog"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	httpErrors "github.com/scul0405/blog-clean-architecture-rest-api/pkg/http_errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"net/http"
)

const (
	basePrefix    = "blog-api"
	cacheDuration = 3600
)

type blogUseCase struct {
	cfg       *config.Config
	blogRepo  blog.Repository
	redisRepo blog.RedisRepository
	logger    logger.Logger
}

func NewBlogUseCase(cfg *config.Config, blogRepo blog.Repository, redisRepo blog.RedisRepository, logger logger.Logger) blog.UseCase {
	return &blogUseCase{cfg: cfg, blogRepo: blogRepo, redisRepo: redisRepo, logger: logger}
}

func (u *blogUseCase) Create(ctx context.Context, blog *models.Blog) (*models.BlogBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogUC.Create")
	defer span.Finish()

	userUID, err := utils.GetUserUIDFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "blogUC.Create.GetUserUIDFromCtx"))
	}

	blog.AuthorID = userUID
	createdBlog, err := u.blogRepo.Create(ctx, blog)
	if err != nil {
		return nil, err
	}

	return createdBlog, nil
}

func (u *blogUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.BlogBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogUC.GetByID")
	defer span.Finish()

	blogCached, err := u.redisRepo.GetBlogByIDCtx(ctx, u.generateBlogKey(id.String()))
	if err != nil {
		u.logger.Errorf("blogUC.GetByID: GetBlogByIDCtx: %v", err)
	}

	if blogCached != nil {
		return blogCached, nil
	}

	blog, err := u.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetBlogCtx(ctx, u.generateBlogKey(blog.BlogID.String()), cacheDuration, blog); err != nil {
		u.logger.Errorf("blogUC.GetByID: SetBlogCtx: %v", err)
	}

	return blog, nil
}

func (u *blogUseCase) Update(ctx context.Context, blog *models.BlogBase) (*models.BlogBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogUC.Update")
	defer span.Finish()

	blogByID, err := u.blogRepo.GetByID(ctx, blog.BlogID)
	if err != nil {
		return nil, err
	}

	if err = utils.ValidateIsOwner(ctx, blogByID.AuthorID.String(), u.logger); err != nil {
		return nil, httpErrors.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "blogUC.Update.ValidateIsOwner"))
	}

	updatedBlog, err := u.blogRepo.Update(ctx, blog)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.DeleteBlogCtx(ctx, u.generateBlogKey(blog.BlogID.String())); err != nil {
		u.logger.Errorf("blogUC.Update.DeleteBlogCtx: %v", err)
	}

	return updatedBlog, nil
}

func (u *blogUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogUC.Delete")
	defer span.Finish()

	blogByID, err := u.blogRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err = utils.ValidateIsOwner(ctx, blogByID.AuthorID.String(), u.logger); err != nil {
		return httpErrors.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "blogUC.Delete.ValidateIsOwner"))
	}

	if err = u.blogRepo.Delete(ctx, id); err != nil {
		return err
	}

	if err = u.redisRepo.DeleteBlogCtx(ctx, u.generateBlogKey(id.String())); err != nil {
		u.logger.Errorf("blogUC.Delete.DeleteBlogCtx: %v", err)
	}

	return nil
}

func (u *blogUseCase) List(ctx context.Context, pq *utils.PaginationQuery) (*models.BlogsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogUC.List")
	defer span.Finish()

	return u.blogRepo.List(ctx, pq)
}

func (u *blogUseCase) generateBlogKey(blogID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, blogID)
}
