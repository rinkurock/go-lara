package middlewares

import (
	"app/helpers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ApplyMiddleware(e *echo.Echo) {
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &helpers.CustomEchoContext{Context: c}
			return h(cc)
		}
	})

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())

	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
}
