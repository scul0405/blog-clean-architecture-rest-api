package http

import (
	"github.com/labstack/echo/v4"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/blog"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/middleware"
)

func MapBlogRoutes(blogGroup *echo.Group, h blog.Handlers, mw *middleware.MiddlewareManager) {
	blogGroup.POST("", h.Create(), mw.AuthPASETOMiddleware)
	blogGroup.GET("/:blog_id", h.GetByID())
	blogGroup.PATCH("/:blog_id", h.Update(), mw.AuthPASETOMiddleware)
	blogGroup.DELETE("/:blog_id", h.Delete(), mw.AuthPASETOMiddleware)
	blogGroup.GET("", h.List())
}
