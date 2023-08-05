package auth

import "github.com/labstack/echo"

type Handlers interface {
	Register() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Login() echo.HandlerFunc
}
