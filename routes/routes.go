package routes

import (
	"app/controllers/auth"
	"app/controllers/sample"

	"github.com/labstack/echo/v4"
)

func DefineRoutes(e *echo.Echo) {
	v1 := e.Group("/v1")

	s := v1.Group("/sample")
	s.GET("", sample.SampleGet)
	s.POST("", sample.SamplePost)
	s.PATCH("/:id", sample.SamplePatch)
	s.PUT("/:id", sample.SamplePut)
	s.DELETE("/:id", sample.SampleDelete)

	v1.GET("/hello/world", auth.HelloWorld)
}
