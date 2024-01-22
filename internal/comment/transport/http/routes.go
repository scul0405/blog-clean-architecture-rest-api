package http

import (
	"github.com/labstack/echo/v4"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/comment"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/middleware"
)

func MapCommentRoutes(commentGroup *echo.Group, h comment.Handlers, mw *middleware.MiddlewareManager) {
	commentGroup.POST("", h.Create(), mw.AuthPASETOMiddleware)
	commentGroup.GET("/:comment_id", h.GetByID())
	commentGroup.PATCH("/:comment_id", h.Update(), mw.AuthPASETOMiddleware)
	commentGroup.DELETE("/:comment_id", h.Delete(), mw.AuthPASETOMiddleware)
	commentGroup.GET("", h.List())
	commentGroup.PATCH("/:comment_id/like", h.Like(), mw.AuthPASETOMiddleware)
	commentGroup.PATCH("/:comment_id/dislike", h.Dislike(), mw.AuthPASETOMiddleware)
}
