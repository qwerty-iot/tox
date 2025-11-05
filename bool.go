package tox

import (
	"strconv"
)

func IfBool[T any](v bool, ifTrue T, ifFalse T) T {
	if v {
		return ifTrue
	}
	return ifFalse
}

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

func ToBoolPtr(v interface{}) *bool {
	if v == nil {
		return nil
	}
	ret := ToBool(v)
	return &ret
}

func TriBool(b *bool) bool {
	if b != nil {
		return *b
	} else {
		return false
	}
}
