package routes

import (
	"app/controllers/auth"

	"github.com/labstack/echo/v4"
)

func DefineRoutes(e *echo.Echo) {
	v1 := e.Group("/v1")
	v1.GET("/hello/world", auth.HelloWorld)
}
