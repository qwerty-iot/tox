package tox

import (
	"strconv"
)

// ToBool converts any data type to a bool, if the conversion fails, it returns false.
func ToBool(v interface{}) bool {
	switch v := v.(type) {
	case bool:
		return v
	case float32:
		if v != 0.0 {
			return true
		}
	case float64:
		if v != 0.0 {
			return true
		}
	case int:
		if v != 0 {
			return true
		}
	case int8:
		if v != 0 {
			return true
		}
	case int16:
		if v != 0 {
			return true
		}
	case int32:
		if v != 0 {
			return true
		}
	case int64:
		if v != 0 {
			return true
		}
	case uint:
		if v != 0 {
			return true
		}
	case uint8:
		if v != 0 {
			return true
		}
	case uint16:
		if v != 0 {
			return true
		}
	case uint32:
		if v != 0 {
			return true
		}
	case uint64:
		if v != 0 {
			return true
		}
	case string:
		i, _ := strconv.ParseBool(v)
		return i
	}
	return false
}
