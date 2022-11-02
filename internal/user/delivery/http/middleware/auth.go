package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func KeyAuth(key []byte) echo.MiddlewareFunc {
	return middleware.JWT(key)
}
