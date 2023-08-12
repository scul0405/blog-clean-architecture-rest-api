package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/blog/mock"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestBlogUseCase_Create(t *testing.T) {
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
	mockBlogRepo := mock.NewMockRepository(ctrl)
	blogUC := NewBlogUseCase(cfg, mockBlogRepo, apiLogger)

	userUID := uuid.New()

	blog := &models.Blog{
		AuthorID: userUID,
		Title:    "Title long text string greater then 20 characters",
		Content:  "Content long text string greater then 20 characters",
	}

	blogBase := &models.BlogBase{
		AuthorID: userUID,
		Title:    "Title long text string greater then 20 characters",
		Content:  "Content long text string greater then 20 characters",
	}

	ctx := context.WithValue(context.Background(), "user_id", userUID.String())
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "blogUC.Create")
	defer span.Finish()

	mockBlogRepo.EXPECT().Create(ctxWithTrace, gomock.Eq(blog)).Return(blogBase, nil)

	createdBlog, err := blogUC.Create(ctx, blog)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, createdBlog)
}

func TestBlogUseCase_GetByID(t *testing.T) {
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
	mockBlogRepo := mock.NewMockRepository(ctrl)
	blogUC := NewBlogUseCase(cfg, mockBlogRepo, apiLogger)

	blogUID := uuid.New()
	userUID := uuid.New()

	blogBase := &models.BlogBase{
		BlogID:   blogUID,
		AuthorID: userUID,
		Title:    "Title long text string greater then 20 characters",
		Content:  "Content long text string greater then 20 characters",
	}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "blogUC.GetByID")
	defer span.Finish()

	mockBlogRepo.EXPECT().GetByID(ctxWithTrace, gomock.Eq(blogUID)).Return(blogBase, nil)

	getByIDBlog, err := blogUC.GetByID(ctx, blogUID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, getByIDBlog)
}

func TestBlogUseCase_Update(t *testing.T) {
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
	mockBlogRepo := mock.NewMockRepository(ctrl)
	blogUC := NewBlogUseCase(cfg, mockBlogRepo, apiLogger)

	blogUID := uuid.New()
	userUID := uuid.New()

	blogBase := &models.BlogBase{
		BlogID:   blogUID,
		AuthorID: userUID,
		Title:    "Title long text string greater then 20 characters",
		Content:  "Content long text string greater then 20 characters",
	}

	ctx := context.WithValue(context.Background(), "user_id", userUID.String())
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "blogUC.Update")
	defer span.Finish()

	mockBlogRepo.EXPECT().GetByID(ctxWithTrace, gomock.Eq(blogUID)).Return(blogBase, nil)
	mockBlogRepo.EXPECT().Update(ctxWithTrace, gomock.Eq(blogBase)).Return(blogBase, nil)

	updatedBlog, err := blogUC.Update(ctx, blogBase)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, updatedBlog)
}

func TestBlogUseCase_Delete(t *testing.T) {
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
	mockBlogRepo := mock.NewMockRepository(ctrl)
	blogUC := NewBlogUseCase(cfg, mockBlogRepo, apiLogger)

	blogUID := uuid.New()
	userUID := uuid.New()

	blogBase := &models.BlogBase{
		BlogID:   blogUID,
		AuthorID: userUID,
		Title:    "Title long text string greater then 20 characters",
		Content:  "Content long text string greater then 20 characters",
	}

	ctx := context.WithValue(context.Background(), "user_id", userUID.String())
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "blogUC.Delete")
	defer span.Finish()

	mockBlogRepo.EXPECT().GetByID(ctxWithTrace, gomock.Eq(blogUID)).Return(blogBase, nil)
	mockBlogRepo.EXPECT().Delete(ctxWithTrace, gomock.Eq(blogUID)).Return(nil)

	err := blogUC.Delete(ctx, blogUID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestBlogUseCase_List(t *testing.T) {
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
	mockBlogRepo := mock.NewMockRepository(ctrl)
	blogUC := NewBlogUseCase(cfg, mockBlogRepo, apiLogger)

	blogUID := uuid.New()
	userUID := uuid.New()

	blogsListMock := &models.BlogsList{
		Blogs: []*models.BlogBase{
			{
				BlogID:   blogUID,
				AuthorID: userUID,
				Title:    "Title long text string greater then 20 characters",
				Content:  "Content long text string greater then 20 characters",
			},
			{
				BlogID:   blogUID,
				AuthorID: userUID,
				Title:    "Title long text string greater then 20 characters",
				Content:  "Content long text string greater then 20 characters",
			},
		},
	}

	pq := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "blogUC.List")
	defer span.Finish()

	mockBlogRepo.EXPECT().List(ctxWithTrace, gomock.Eq(pq)).Return(blogsListMock, nil)

	blogsList, err := blogUC.List(ctx, pq)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, blogsList)
	require.Equal(t, len(blogsList.Blogs), 2)
}
