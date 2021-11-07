package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type CustomEchoContext struct {
	echo.Context
	allFormValuesInMap     map[string]interface{}
	allJsonFormValuesInMap map[string]interface{}
	allQueryParams         map[string]interface{}
}

// func (cc *CustomEchoContext) HasData() bool {
// 	if _, err := cc.FormValuesInMap(); err == nil {
// 		return true
// 	}
// 	return false
// }

func (cc *CustomEchoContext) HasFormKey(key string) bool {
	if all, err := cc.FormValuesInMap(); err == nil {
		if MapValue(key, all) != nil {
			return true
		}
	}
	return false
}

func (cc *CustomEchoContext) HasQueryParamKey(key string) bool {
	if cc.allQueryParams == nil {
		cc.allQueryParams = cc.AllQueryParamsInMap()
	}
	if MapValue(key, cc.allQueryParams) != nil {
		return true
	}

	return false
}

func (cc *CustomEchoContext) HasJsonFormKey(key string) bool {
	if all, err := cc.JsonFormValuesInMap(); err == nil {
		if MapValue(key, all) != nil {
			return true
		}
	}
	return false
}

func (cc *CustomEchoContext) JsonFormKey(key string) interface{} {
	if all, err := cc.JsonFormValuesInMap(); err == nil {
		if result := MapValue(key, all); result != nil {
			return result
		}
	}
	return false
}

func (cc *CustomEchoContext) JsonFormKeyString(key string) string {
	if all, err := cc.JsonFormValuesInMap(); err == nil {
		if result := MapValue(key, all); result != nil {
			if reflect.TypeOf(result).Kind() == reflect.String {
				return result.(string)
			}
		}
	}
	return ""
}

func (cc *CustomEchoContext) HasFormFile(fileName string) bool {
	if _, err := cc.FormFile(fileName); err == nil {
		return true
	}
	return false
}

func (cc *CustomEchoContext) FormFileBase64Encoded(uploadedFileName string) string {
	if fileHeader, err := cc.FormFile(uploadedFileName); err == nil {
		if file, err := fileHeader.Open(); err == nil {
			defer func() { _ = file.Close() }()
			if fileByte, err := ioutil.ReadAll(file); err == nil {
				return base64.StdEncoding.EncodeToString(fileByte)
			} else {
				log.Error("could not read fileHeader: ", uploadedFileName)
			}
		} else {
			log.Error("could not open uploaded file : ", uploadedFileName)
		}
	} else {
		log.Error("could not get fileHeader from request: ", uploadedFileName)
	}
	return ""
}

func (cc *CustomEchoContext) FormValuesInMap() (allFormValues map[string]interface{}, err error) {
	if cc.allFormValuesInMap != nil {
		allFormValues = cc.allFormValuesInMap
	} else {
		f := url.Values{}
		f, err = cc.FormParams()
		allFormValues = map[string]interface{}{}
		for key, value := range f {
			allFormValues[key] = value
		}
	}
	return
}

func (cc *CustomEchoContext) JsonFormValuesInMap() (allJsonFormValues map[string]interface{}, err error) {
	if cc.allJsonFormValuesInMap != nil {
		allJsonFormValues = cc.allJsonFormValuesInMap
	} else {
		//Bind will work if only sent via application-json
		err = cc.Bind(&allJsonFormValues)
		cc.allJsonFormValuesInMap = allJsonFormValues
	}
	return
}

func (cc *CustomEchoContext) JsonFormValuesWithoutKeys(keysToRemove string) (result map[string]interface{}, err error) {
	if jsonFormValuesInMap, err := cc.JsonFormValuesInMap(); err == nil {
		result, err = MapWithoutKeys(jsonFormValuesInMap, keysToRemove)
	}
	return
}

func (cc *CustomEchoContext) FormValuesWithoutKeys(keysToRemove string) (result map[string]interface{}, err error) {
	if formValuesInMap, err := cc.FormValuesInMap(); err == nil {
		result, err = MapWithoutKeys(formValuesInMap, keysToRemove)
	}
	return
}

func (cc *CustomEchoContext) FormValuesWithOnlyKeys(keysToKeep string) (result map[string]interface{}, err error) {
	if formValuesInMap, err := cc.FormValuesInMap(); err == nil {
		result = MapOnlyKeys(formValuesInMap, keysToKeep)
	}
	return
}

func (cc *CustomEchoContext) FillStructWithOnlyFormValues(onlyFormKeys string, structToFill interface{}) error {
	if formValues, err := cc.FormValuesWithOnlyKeys(onlyFormKeys); err == nil {
		return MapToStruct(formValues, &structToFill)
	} else {
		return err
	}
}

func (cc *CustomEchoContext) FillStructWithOutFormValues(keysToRemove string, output interface{}) error {
	if formValues, err := cc.FormValuesWithoutKeys(keysToRemove); err == nil {
		return MapToStruct(formValues, &output)
	} else {
		return err
	}
}

func (cc *CustomEchoContext) FormValueInt(keyName string) int {
	val := cc.FormValue(keyName)
	if number, err := strconv.Atoi(val); err == nil {
		return number
	}
	return 0
}

func (cc *CustomEchoContext) FormValueInt32(keyName string) int32 {
	val := cc.FormValue(keyName)
	if number, err := strconv.ParseInt(val, 10, 32); err == nil {
		return int32(number)
	}
	return 0
}
func MapToStruct(inputMap map[string]interface{}, outputStuct interface{}) error {
	if b, err := json.Marshal(inputMap); err == nil {
		return json.Unmarshal(b, &outputStuct)
	} else {
		return err
	}
}

func (cc *CustomEchoContext) FormValueInt64(keyName string) int64 {
	val := cc.FormValue(keyName)
	if number, err := strconv.ParseInt(val, 10, 64); err == nil {
		return number
	}
	return 0
}

func (cc *CustomEchoContext) FormValueFloat32(keyName string) float32 {
	val := cc.FormValue(keyName)
	if number, err := strconv.ParseFloat(val, 32); err == nil {
		return float32(number)
	}
	return 0
}

func (cc *CustomEchoContext) FormValueFloat64(keyName string) float64 {
	val := cc.FormValue(keyName)
	if number, err := strconv.ParseFloat(val, 64); err == nil {
		return number
	}
	return 0
}

func (cc *CustomEchoContext) FormValueBool(keyName string) bool {
	val := cc.FormValue(keyName)
	if res, err := strconv.ParseBool(val); err == nil {
		return res
	}
	return false
}

func (cc *CustomEchoContext) FillStructWithFormValues(output interface{}) error {
	return cc.Bind(&output)
}

func (cc *CustomEchoContext) BindWeakJson(output interface{}, isZeroFields ...bool) (err error) {
	zeroFields := false
	if len(isZeroFields) > 0 {
		zeroFields = isZeroFields[0]
	}

	if data, err := cc.JsonFormValuesInMap(); err == nil && len(data) > 0 {
		if err = MapToStructWeak(data, &output, zeroFields); err != nil {
			log.Error(err)
		}
	}

	return
}

func (cc *CustomEchoContext) BindWeakQuery(output interface{}, isZeroFields ...bool) (err error) {
	zeroFields := false
	if len(isZeroFields) > 0 {
		zeroFields = isZeroFields[0]
	}

	if len(cc.AllQueryParamsInMap()) > 0 {
		err = MapToStructWeak(cc.AllQueryParamsInMap(), &output, zeroFields)
	}

	if err != nil {
		log.Error(err)
	}

	return
}

func (cc *CustomEchoContext) BindSafe(output interface{}) (err error) {
	if cc.Request().ContentLength != 0 && len(cc.allJsonFormValuesInMap) < 1 && len(cc.allFormValuesInMap) < 1 {
		err = cc.Bind(output)
		contentType := cc.Request().Header.Get(echo.HeaderContentType)
		if values, err := StructToMap(&output); err == nil {
			if strings.HasPrefix(contentType, echo.MIMEApplicationJSON) {
				cc.allJsonFormValuesInMap = values
			} else if strings.HasPrefix(contentType, echo.MIMEApplicationForm) {
				cc.allFormValuesInMap = values
			}
		}
	} else {
		kind := reflect.TypeOf(output).Kind()

		if len(cc.allJsonFormValuesInMap) > 0 {
			if kind == reflect.Map {
				output = cc.allJsonFormValuesInMap
				return nil
			} else if kind == reflect.Struct || kind == reflect.Ptr {
				return MapToStruct(cc.allJsonFormValuesInMap, output)
			}
		} else if len(cc.allFormValuesInMap) > 0 {
			if kind == reflect.Map {
				output = cc.allFormValuesInMap
				return nil
			} else if kind == reflect.Struct || kind == reflect.Ptr {
				return MapToStruct(cc.allFormValuesInMap, output)
			}
		} else {
			return nil
			//TODO return error if no body is found
		}
	}
	return
}

func (cc *CustomEchoContext) AllQueryParamsInMap() map[string]interface{} {
	if cc.allQueryParams == nil {
		cc.allQueryParams = map[string]interface{}{}
		for key, values := range cc.QueryParams() {
			if len(values) > 0 {
				cc.allQueryParams[key] = values[0]
			}
		}
	}

	return cc.allQueryParams
}

func (cc *CustomEchoContext) FormOrQueryParamStringValue(key string) string {
	log.Info(cc.allQueryParams)
	log.Info(cc.QueryParam(key))
	if cc.HasFormKey(key) {
		return cc.FormValue(key)
	}
	if cc.HasQueryParamKey(key) {
		return cc.QueryParam(key)
	}
	return ""
}

/*************** Helper for Custom Echo Context ******************/

func MapValue(key string, mapFrom map[string]interface{}) interface{} {
	if value, ok := mapFrom[key]; ok {
		return value
	}
	return nil
}

func MapWithoutKeys(sourceMap map[string]interface{}, keysToRemove string) (map[string]interface{}, error) {
	if keysToRemove != "" {
		keysToRemove := strings.Split(keysToRemove, ",")
		for _, value := range keysToRemove {
			MapDelKey(value, sourceMap)
		}
		return sourceMap, nil
	} else {
		return sourceMap, nil
	}
}
func MapDelKey(keyToDel string, mapFrom map[string]interface{}) {
	if _, ok := mapFrom[keyToDel]; ok {
		delete(mapFrom, keyToDel)
	}
}

func MapOnlyKeys(sourceMap map[string]interface{}, onlyKeys string) map[string]interface{} {
	resultMap := map[string]interface{}{}
	if len(onlyKeys) > 0 {
		keysToKeep := strings.Split(onlyKeys, ",")
		for _, key := range keysToKeep {
			val := MapValue(key, sourceMap)
			if val != nil {
				resultMap[key] = val
			}
		}
	}
	return resultMap
}
func StructToMap(st interface{}) (map[string]interface{}, error) {
	if b, err := json.Marshal(st); err == nil {
		var result map[string]interface{}
		if err := json.Unmarshal(b, &result); err == nil {
			return result, nil
		}
	}
	return nil, errors.New("could not convert to map")
}
