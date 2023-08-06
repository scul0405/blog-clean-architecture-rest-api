package utils

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
)

// GetRequestID get the request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// GetIPAddress get the ip address from echo context
func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

// ReqIDCtxKey is a key used for the Request ID in context
type ReqIDCtxKey struct{}

// GetRequestCtx get context with request id
func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
}

// ReadRequest read request body and validate
func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request().Context(), request)
}

func LogResponseError(ctx echo.Context, logger logger.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	)
}
