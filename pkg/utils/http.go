package utils

import "github.com/labstack/echo"

// GetRequestID get the request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}
