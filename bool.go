package tox

import (
	"strconv"
)

func ToBool(v interface{}) bool {
	switch v := v.(type) {
	case bool:
		return v
	case float32:
		if v != 0.0 {
			return true
		} else {
			return false
		}
	case float64:
		if v != 0.0 {
			return true
		} else {
			return false
		}
	case int:
		if v != 0 {
			return true
		} else {
			return false
		}
	case int8:
		if v != 0 {
			return true
		} else {
			return false
		}
	case int16:
		if v != 0 {
			return true
		} else {
			return false
		}
	case int32:
		if v != 0 {
			return true
		} else {
			return false
		}
	case int64:
		if v != 0 {
			return true
		} else {
			return false
		}
	case uint:
		if v != 0 {
			return true
		} else {
			return false
		}
	case uint8:
		if v != 0 {
			return true
		} else {
			return false
		}
	case uint16:
		if v != 0 {
			return true
		} else {
			return false
		}
	case uint32:
		if v != 0 {
			return true
		} else {
			return false
		}
	case uint64:
		if v != 0 {
			return true
		} else {
			return false
		}
	case string:
		i, _ := strconv.ParseBool(v)
		return i
	default:
		return false
	}
}
