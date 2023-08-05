package http

import (
	"github.com/labstack/echo"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/auth"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers) {
	authGroup.POST("/register", h.Register())
	authGroup.GET("/:id", h.GetByID())
	authGroup.POST("/login", h.Login())
}
