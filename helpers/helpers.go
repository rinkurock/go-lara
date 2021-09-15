package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/mitchellh/mapstructure"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func ToString(itemToConvert interface{}) string {
	switch reflect.TypeOf(itemToConvert).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", itemToConvert)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.2f", itemToConvert)
	case reflect.Bool:
		return strconv.FormatBool(itemToConvert.(bool))
	case reflect.Array, reflect.Struct:
		return fmt.Sprintf("%v", itemToConvert)
	default:
		return ""
	}
}
func ByteArrayToMap(bytes []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err == nil {
		return result, nil
	}
	return nil, errors.New("could not convert to map")
}
func MapToStructWeak(inputMap map[string]interface{}, outputStruct interface{}, isZeroFields bool) (err error) {
	var md mapstructure.Metadata
	config := &mapstructure.DecoderConfig{
		Metadata:         &md,
		Result:           &outputStruct,
		ZeroFields:       isZeroFields,
		WeaklyTypedInput: true,
		TagName:          "json",
	}
	if decoder, err := mapstructure.NewDecoder(config); err == nil {
		err = decoder.Decode(&inputMap)
		log.Debug(fmt.Sprintf("keys: %#v", md.Keys))
		log.Debug(fmt.Sprintf("Unused: %#v", md.Unused))
	}
	return
}
