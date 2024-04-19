package tox

import (
	"reflect"
	"strconv"
)

// ToInt converts any data type to a bool, if the conversion fails, it returns 0.
func ToInt(v interface{}) int {
	switch v := v.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		i, _ := strconv.Atoi(v)
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

func ToIntPtr(v interface{}) *int {
	if v == nil {
		return nil
	}
	ret := ToInt(v)
	return &ret
}

func ToIntArray(v interface{}) []int {
	switch v := v.(type) {
	case nil:
		return nil
	case int:
		return []int{v}
	case []int:
		return v
	case []byte:
		var ret = make([]int, len(v))
		for ii, vv := range v {
			ret[ii] = int(vv)
		}
		return ret
	case []any:
		var ret = make([]int, len(v))
		for ii, vv := range v {
			ret[ii] = ToInt(vv)
		}
		return ret
	default:
		aVal := reflect.ValueOf(v)
		if aVal.Kind() == reflect.Array || aVal.Kind() == reflect.Slice {
			var ret = make([]int, aVal.Len())
			for i := 0; i < aVal.Len(); i++ {
				ret[i] = ToInt(aVal.Index(i).Interface())
			}
			return ret
		} else {
			return []int{ToInt(v)}
		}
	}
}
