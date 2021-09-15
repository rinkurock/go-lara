package helpers

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
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
