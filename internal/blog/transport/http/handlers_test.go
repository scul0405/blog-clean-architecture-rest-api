package http

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/blog/mock"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/converter"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestBlogHandlers_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogUC := mock.NewMockUseCase(ctrl)

	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	blogHandlers := NewBlogHandlers(cfg, mockBlogUC, apiLogger)

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

	buf, err := converter.ToBytesBuffer(blog)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, buf)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/api/v1/blogs", strings.NewReader(buf.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctxWithValue := context.WithValue(context.Background(), "user_id", userUID)
	req.WithContext(ctxWithValue)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "blogHandlers.Create")
	defer span.Finish()

	mockBlogUC.EXPECT().Create(ctxWithTrace, gomock.Eq(blog)).Return(blogBase, nil)

	handlerFunc := blogHandlers.Create()
	err = handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestBlogHandlers_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogUC := mock.NewMockUseCase(ctrl)

	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	blogHandlers := NewBlogHandlers(cfg, mockBlogUC, apiLogger)

	blogUID := uuid.New()
	userUID := uuid.New()
	blogBase := &models.BlogBase{
		BlogID:   blogUID,
		AuthorID: userUID,
		Title:    "Title long text string greater then 20 characters",
		Content:  "Content long text string greater then 20 characters",
	}

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/blogs/:blog_id", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("blog_id")
	c.SetParamValues(blogUID.String())
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "blogHandlers.GetByID")
	defer span.Finish()

	mockBlogUC.EXPECT().GetByID(ctxWithTrace, gomock.Eq(blogUID)).Return(blogBase, nil)

	handlerFunc := blogHandlers.GetByID()
	err := handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestBlogHandlers_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogUC := mock.NewMockUseCase(ctrl)

	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	blogHandlers := NewBlogHandlers(cfg, mockBlogUC, apiLogger)

	blogUID := uuid.New()
	userUID := uuid.New()
	blogBase := &models.BlogBase{
		BlogID:   blogUID,
		AuthorID: userUID,
		Title:    "Title long text string greater then 20 characters",
		Content:  "Content long text string greater then 20 characters",
	}

	buf, err := converter.ToBytesBuffer(blogBase)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, buf)

	e := echo.New()
	req := httptest.NewRequest(echo.PATCH, "/api/v1/blogs/:blog_id", strings.NewReader(buf.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctxWithValue := context.WithValue(context.Background(), "user_id", userUID)
	req.WithContext(ctxWithValue)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("blog_id")
	c.SetParamValues(blogUID.String())
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "blogHandlers.Update")
	defer span.Finish()

	mockBlogUC.EXPECT().Update(ctxWithTrace, gomock.Eq(blogBase)).Return(blogBase, nil)

	handlerFunc := blogHandlers.Update()
	err = handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestBlogHandlers_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogUC := mock.NewMockUseCase(ctrl)

	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	blogHandlers := NewBlogHandlers(cfg, mockBlogUC, apiLogger)

	blogUID := uuid.New()
	userUID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/api/v1/blogs/:blog_id", nil)
	ctxWithValue := context.WithValue(context.Background(), "user_id", userUID)
	req.WithContext(ctxWithValue)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("blog_id")
	c.SetParamValues(blogUID.String())
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "blogHandlers.Delete")
	defer span.Finish()

	mockBlogUC.EXPECT().Delete(ctxWithTrace, gomock.Eq(blogUID)).Return(nil)

	handlerFunc := blogHandlers.Delete()
	err := handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestBlogHandlers_List(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogUC := mock.NewMockUseCase(ctrl)

	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	blogHandlers := NewBlogHandlers(cfg, mockBlogUC, apiLogger)

	pq := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	blogUID := uuid.New()
	userUID := uuid.New()
	blogList := &models.BlogsList{
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

	e := echo.New()
	q := make(url.Values)
	q.Set("page", "1")
	q.Set("size", "10")

	req := httptest.NewRequest(echo.GET, "/api/v1/blogs?"+q.Encode(), nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "blogHandlers.List")
	defer span.Finish()

	mockBlogUC.EXPECT().List(ctxWithTrace, gomock.Eq(pq)).Return(blogList, nil)

	handlerFunc := blogHandlers.List()
	err := handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}
