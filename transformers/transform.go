package transform

import (
	h "app/helpers"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func SuccessResponse(c *h.CustomEchoContext, data map[string]interface{}) error {
	return c.JSON(http.StatusOK, data)
}

func UnauthorizedResponse(c *h.CustomEchoContext) error {
	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
		"error":             "access_denied",
		"error_description": "The resource owner or authorization server denied the request.",
	})
}

func BodyParsingErrorResponse(c *h.CustomEchoContext, err error) error {
	log.Errorln("bodyParsingErrorResponse", err.Error())
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error": "request body parsing error",
	})
}
func NotFoundResponse(c *h.CustomEchoContext, message string) error {
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"error": message,
	})
}
