package server

import (
	"github.com/labstack/echo"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"net/http"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	v1 := e.Group("/api/v1")

	health := v1.Group("/health")

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})
	return nil
}
