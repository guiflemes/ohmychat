package utils

import (
	"reflect"
	"strconv"
)

func Parse(value any) string {
	t := reflect.TypeOf(value)

	if t.Kind() == reflect.Slice {
		return Slice(value.([]any), t)
	}
	return Primitive(value, t)
}

func Primitive(value any, t reflect.Type) string {

	switch t.Kind() {
	case reflect.Bool:
		return Boolean(value.(bool))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Integer(value.(int))

	case reflect.Float32, reflect.Float64:
		return Float64(value.(float64))

	case reflect.Map:
		return Mapping(value.(map[string]any))

	default:
		return value.(string)
	}

}

func Mapping(data map[string]any) string {
	filteredMap := ""
	for key, value := range data {
		switch val := value.(type) {
		case string:
			filteredMap += key + ": " + val + "\n"
		case int, int8, int16, int64:
			filteredMap += key + ": " + Integer(value.(int)) + "\n"
		case float32, float64:
			filteredMap += key + ": " + Float64(value.(float64)) + "\n"
		case bool:
			filteredMap += key + ": " + Boolean(value.(bool)) + "\n"
		default:
			// Skip slices, maps, and other types
		}
	}
	return filteredMap
}

func Slice(value []any, t reflect.Type) string {
	v := ""
	for _, item := range value {
		str := Primitive(item, t.Elem())
		v += str + ", "
	}
	return v
}

func Integer(value int) string { return strconv.Itoa(value) }

func Float64(value float64) string { return strconv.FormatFloat(value, 'f', -1, 64) }

func Boolean(value bool) string { return strconv.FormatBool(value) }
