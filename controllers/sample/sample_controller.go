package sample

import (
	h "app/helpers"
	"app/models"
	repo "app/repositories"
	t "app/transformers"
	validate "app/validation"
	"net/http"

	// "github.com/davecgh/go-spew/spew"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func SampleGet(echoContext echo.Context) error {
	c := echoContext.(*h.CustomEchoContext)
	res := repo.SampleGetResponse()
	return c.JSON(http.StatusOK, res)
}

func SamplePost(echoContext echo.Context) error {
	c := echoContext.(*h.CustomEchoContext)
	var req models.SamplePostReq
	err := c.BindWeakJson(&req)
	if err != nil {
		log.Errorln("request body parsing error  on SamplePost", err.Error())
		return t.BodyParsingErrorResponse(c, err)
	}
	v := _validateSamplePost(c, req)
	if v.Status {
		return c.JSON(v.StatusCode, v.Response)
	}
	res := repo.GetPostResponse(req)
	return c.JSON(http.StatusOK, res)
}

func SamplePatch(echoContext echo.Context) error {
	c := echoContext.(*h.CustomEchoContext)
	var req models.SamplePatchReq
	id := h.ToInt(c.Param("id"))
	if id != 99 {
		return t.NotFoundResponse(c, "Profile not found!")
	}
	err := c.Bind(&req)
	if err != nil {
		log.Errorln("request body parsing error  on SamplePost", err.Error())
		return t.BodyParsingErrorResponse(c, err)
	}
	v := _validateSamplePatch(c, req)
	if v.Status {
		return c.JSON(v.StatusCode, v.Response)
	}

	res := repo.PatchResponse(req)
	return c.JSON(http.StatusOK, res)
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
	if !c.HasJsonFormKey("term_accepted") || !data.TermAccepted{
		m["term_accepted"] = "valid term_accepted required"
		v.Status = true
	}
	v.Response = m
	return v
}

func _validateSamplePatch(c *h.CustomEchoContext, data models.SamplePatchReq) validate.ValidationError {
	v := validate.NewValidationError()
	m := make(map[string]interface{})
	if c.HasFormKey("name") {
		if len(data.Name) > 255 || len(data.Name) < 1 {
			m["name"] = "valid name required"
			v.Status = true
		}
	}
	if c.HasFormKey("number") {
		if len(data.Number) > 20 || len(data.Number) < 10 {
			m["number"] = "valid number required"
			v.Status = true
		}
	}
	if c.HasFormKey("country_id") {
		if data.CountryId != 1 {
			m["country_id"] = "valid country_id required"
			v.Status = true
		}
	}
	if c.HasFormKey("code") {
		if data.Code < 100000 || data.Code > 999999 {
			m["code"] = "valid code required"
			v.Status = true
		}
	}
	if c.HasFormKey("term_accepted") {
		if !data.TermAccepted{
			m["term_accepted"] = "valid term_accepted required"
			v.Status = true
		}
	}
	v.Response = m
	return v
}
