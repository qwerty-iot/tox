package tox

import (
	"strconv"
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
	switch i.(type) {
	case []byte:
		return true
	case string:
		return true
	default:
		return false
	}
}
