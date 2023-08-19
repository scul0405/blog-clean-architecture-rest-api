package http

import (
	"github.com/labstack/echo/v4"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/auth"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/middleware"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.GET("/:id", h.GetByID())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/:id/avatar", h.UploadAvatar())
}
