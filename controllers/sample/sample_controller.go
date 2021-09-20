package sample

import (
	h "app/helpers"
	"app/models"
	repo "app/repositories"
	t "app/transformers"
	validate "app/validation"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func SampleGet(echoContext echo.Context) error {
	c := echoContext.(*h.CustomEchoContext)
	var data map[string]interface{}
	data["hello"] = "world"
	return t.SuccessResponse(c, data)
}

func SamplePost(echoContext echo.Context) error {
	c := echoContext.(*h.CustomEchoContext)
	var req models.SamplePostReq
	err := c.BindWeakJson(&req)
	if err != nil {
		log.Errorln("request body parsing error  on SamplePost", err.Error())
		return t.BodyParsingErrorResponse(c, err)
	}
	spew.Dump(req)
	v := _validateSamplePost(c, req)
	if v.Status {
		return c.JSON(v.StatusCode, v.Response)
	}

	res := repo.GetPostResponse(req)
	return c.JSON(http.StatusOK, res)
}

func _validateSamplePost(c *h.CustomEchoContext, data models.SamplePostReq) validate.ValidationError {
	v := validate.NewValidationError()
	m := make(map[string]interface{})
	if !c.HasJsonFormKey("name") || len(data.Name) > 255 || len(data.Name) < 1 {
		m["name"] = "valid name required"
		v.Status = true
	}
	if !c.HasJsonFormKey("number") || len(data.Number) > 20 || len(data.Number) < 10 {
		m["number"] = "valid number required"
		v.Status = true
	}
	if !c.HasJsonFormKey("country_id") || data.CountryId != 1 {
		m["country_id"] = "valid country_id required"
		v.Status = true
	}
	if !c.HasJsonFormKey("code") || data.Code < 100000 || data.Code > 999999 {
		m["code"] = "valid code required"
		v.Status = true
	}
	if !c.HasJsonFormKey("code") || data.TermAccepted != true {
		m["term_accepted"] = "valid term_accepted required"
		v.Status = true
	}
	v.Response = m
	return v
}

func SamplePatch(echoContext echo.Context) error {
	c := echoContext

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Hello": "world",
	})
}

func SamplePut(echoContext echo.Context) error {
	c := echoContext

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Hello": "world",
	})
}

func SampleDelete(echoContext echo.Context) error {
	c := echoContext

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Hello": "world",
	})
}
