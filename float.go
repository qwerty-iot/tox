package tox

import (
	"math"
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
		i, _ := strconv.ParseFloat(v, 64)
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
