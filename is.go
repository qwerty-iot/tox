package tox

import (
	"reflect"
	"strconv"
	"unicode"
)

func IsNumber(i any) bool {
	switch v := i.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	case string:
		if _, err := strconv.Atoi(v); err != nil {
			return false
		} else {
			return true
		}
	default:
		return false
	}
}

func IsString(i any) bool {
	switch v := i.(type) {
	case []byte:
		for _, b := range string(v) {
			if !unicode.IsPrint(b) {
				return false
			}
		}
		return true
	case string:
		return true
	default:
		return false
	}
}

func IsArray(i any) bool {
	if reflect.TypeOf(i).Kind() == reflect.Array || reflect.TypeOf(i).Kind() == reflect.Slice {
		return true
	} else if reflect.TypeOf(i).Kind() == reflect.Ptr && (reflect.TypeOf(i).Elem().Kind() == reflect.Array || reflect.TypeOf(i).Elem().Kind() == reflect.Slice) {
		return true
	} else {
		return false
	}
}
