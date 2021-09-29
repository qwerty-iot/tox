package tox

import "strconv"

// ToInt converts any data type to a int64, if the conversion fails, it returns 0.
func ToInt64(v interface{}) int64 {
	switch v := v.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func ToInt64Ptr(v interface{}) *int64 {
	if v == nil {
		return nil
	}
	ret := ToInt64(v)
	return &ret
}
