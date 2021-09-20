package routes

import (
	"app/controllers/auth"
	"app/controllers/sample"

	"github.com/labstack/echo/v4"
)

func DefineRoutes(e *echo.Echo) {
	v1 := e.Group("/v1")

	s := v1.Group("/sample")
	s.GET("/get", sample.SampleGet)
	s.POST("/post", sample.SamplePost)
	s.PUT("/post", sample.SamplePut)
	s.PATCH("/post", sample.SamplePatch)
	s.DELETE("/post", sample.SampleDelete)

	v1.GET("/hello/world", auth.HelloWorld)
}
