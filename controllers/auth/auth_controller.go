package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HelloWorld(echoContext echo.Context) error {
	c := echoContext

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Hello": "world",
	})
}
