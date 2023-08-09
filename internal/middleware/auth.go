package middleware

import (
	"context"
	"github.com/labstack/echo/v4"
	httpErrors "github.com/scul0405/blog-clean-architecture-rest-api/pkg/http_errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/paseto"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func (mw *MiddlewareManager) AuthPASETOMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearerHeader := c.Request().Header.Get("Authorization")
		if bearerHeader == "" {
			mw.logger.Error("auth middleware", zap.String("bearerHeader", "bearerHeader = \"\""))
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		headerParts := strings.Split(bearerHeader, " ")
		if len(headerParts) != 2 {
			mw.logger.Error("auth middleware", zap.String("headerParts", "len(headerParts) != 2"))
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		tokenString := headerParts[1]
		payload, err := paseto.VerifyPASETOToken(tokenString, mw.cfg)
		if err != nil {
			mw.logger.Error("auth middleware", zap.String("verifyPASETO", err.Error()))
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		c.Set("user_id", payload.ID)

		ctx := context.WithValue(c.Request().Context(), "user_id", payload.ID)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
