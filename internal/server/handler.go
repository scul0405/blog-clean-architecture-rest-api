package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	authRepo "github.com/scul0405/blog-clean-architecture-rest-api/internal/auth/repository"
	authHttp "github.com/scul0405/blog-clean-architecture-rest-api/internal/auth/transport/http"
	authUC "github.com/scul0405/blog-clean-architecture-rest-api/internal/auth/usecase"
	blogRepo "github.com/scul0405/blog-clean-architecture-rest-api/internal/blog/repository"
	blogHttp "github.com/scul0405/blog-clean-architecture-rest-api/internal/blog/transport/http"
	blogUC "github.com/scul0405/blog-clean-architecture-rest-api/internal/blog/usecase"
	commentRepo "github.com/scul0405/blog-clean-architecture-rest-api/internal/comment/repository"
	commentHttp "github.com/scul0405/blog-clean-architecture-rest-api/internal/comment/transport/http"
	commentUC "github.com/scul0405/blog-clean-architecture-rest-api/internal/comment/usecase"
	apiMiddleware "github.com/scul0405/blog-clean-architecture-rest-api/internal/middleware"
	"strings"

	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"net/http"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init	repositories
	authRepo := authRepo.NewAuthRepository(s.db)
	blogRepo := blogRepo.NewBlogRepository(s.db)
	commentRepo := commentRepo.NewCommentRepository(s.db)

	// Init use cases
	authUC := authUC.NewAuthUseCase(s.cfg, authRepo, s.logger)
	blogUC := blogUC.NewBlogUseCase(s.cfg, blogRepo, s.logger)
	commentUC := commentUC.NewCommentUseCase(s.cfg, commentRepo, s.logger)

	// Init handlers
	authHandler := authHttp.NewAuthHandlers(s.cfg, authUC, s.logger)
	blogHandler := blogHttp.NewBlogHandlers(s.cfg, blogUC, s.logger)
	commentHandler := commentHttp.NewCommentHandlers(s.cfg, commentUC, s.logger)

	// Group routes
	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	blogGroup := v1.Group("/blogs")
	commentGroup := v1.Group("/comments")

	// API middleware
	mw := apiMiddleware.NewMiddlewareManager(authUC, s.cfg, s.logger)
	e.Use(mw.RequestLoggerMiddleware)

	// echo middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))

	e.Use(middleware.RequestID())

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger") // TODO: Add swagger
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	// Map routes
	authHttp.MapAuthRoutes(authGroup, authHandler, mw)
	blogHttp.MapBlogRoutes(blogGroup, blogHandler, mw)
	commentHttp.MapCommentRoutes(commentGroup, commentHandler, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})
	return nil
}
