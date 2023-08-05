package server

import (
	"github.com/labstack/echo"
	authRepo "github.com/scul0405/blog-clean-architecture-rest-api/internal/auth/repository"
	authHttp "github.com/scul0405/blog-clean-architecture-rest-api/internal/auth/transport/http"
	authUC "github.com/scul0405/blog-clean-architecture-rest-api/internal/auth/usecase"

	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"net/http"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init	repositories
	authRepo := authRepo.NewAuthRepository(s.db)

	// Init use cases
	authUC := authUC.NewAuthUseCase(s.cfg, authRepo, s.logger)

	// Init handlers
	authHandler := authHttp.NewAuthHandlers(s.cfg, authUC, s.logger)

	// Group routes
	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")

	// Map routes
	authHttp.MapAuthRoutes(authGroup, authHandler)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})
	return nil
}
