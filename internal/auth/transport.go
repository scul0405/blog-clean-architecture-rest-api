package auth

import "github.com/labstack/echo/v4"

type Handlers interface {
	Register() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Login() echo.HandlerFunc
	UploadAvatar() echo.HandlerFunc
}
