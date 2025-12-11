package tox

import (
	"reflect"
	"strconv"
	"time"
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

func IsPrimitive(i any) bool {
	if i == nil {
		return false
	}

	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	// Numeric primitives
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Bool,
		reflect.String:
		return true

	case reflect.Struct:
		// Allow time.Time as primitive without importing time
		if t.PkgPath() == "time" && t.Name() == "Time" {
			return true
		}
		return false
	default:
		return false
	}
}

func IsTime(i any) bool {
	if reflect.TypeOf(i).PkgPath() == "time" && reflect.TypeOf(i).Name() == "Time" {
		return true
	} else if reflect.TypeOf(i).Kind() == reflect.Ptr && reflect.TypeOf(i).Elem().PkgPath() == "time" && reflect.TypeOf(i).Elem().Name() == "Time" {
		return true
	} else if reflect.TypeOf(i).Kind() == reflect.Interface && reflect.TypeOf(i).Elem().PkgPath() == "time" && reflect.TypeOf(i).Elem().Name() == "Time" {
		return true
	} else if IsString(i) {
		t, err := time.Parse(time.RFC3339Nano, ToString(i))
		return err == nil && !t.IsZero()
	}
	return false
}
