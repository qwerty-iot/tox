package tox

import (
	"math"
	"reflect"
	"strconv"
)

// ToFloat64 converts any data type to a float64, if the conversion fails, it returns NaN.
func ToFloat64(v interface{}) float64 {
	switch v := v.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return math.NaN()
		}
		return i
	case bool:
		if v {
			return float64(1)
		}
		return float64(0)
	default:
		return math.NaN()
	}
}

func ToFloat64Ptr(v interface{}) *float64 {
	if v == nil {
		return nil
	}
	ret := ToFloat64(v)
	return &ret
}

func ToFloat64Array(v interface{}) []float64 {
	switch v := v.(type) {
	case nil:
		return nil
	case float64:
		return []float64{v}
	case []float64:
		return v
	case []byte:
		var ret = make([]float64, len(v))
		for ii, vv := range v {
			ret[ii] = float64(vv)
		}
		return ret
	case []any:
		var ret = make([]float64, len(v))
		for ii, vv := range v {
			ret[ii] = ToFloat64(vv)
		}
		return ret
	default:
		aVal := reflect.ValueOf(v)
		if aVal.Kind() == reflect.Array || aVal.Kind() == reflect.Slice {
			var ret = make([]float64, aVal.Len())
			for i := 0; i < aVal.Len(); i++ {
				ret[i] = ToFloat64(aVal.Index(i).Interface())
			}
			return ret
		} else {
			return []float64{ToFloat64(v)}
		}
	}
}
