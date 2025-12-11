package tox

import (
    "math"
    "strconv"
)

func ToPtr[T any](a T) *T {
	p := new(T)
	*p = a
	return p
}

func ToNumber(v any) any {
    minInt, maxInt := platformIntLimits()
    switch val := v.(type) {
    case float32:
        f := float64(val)
        if math.Trunc(f) == f { // integral value
            if f >= float64(minInt) && f <= float64(maxInt) {
                return int(f)
            }
            if f >= float64(math.MinInt64) && f <= float64(math.MaxInt64) {
                return int64(f)
            }
        }
        return f
    case float64:
        f := val
        if math.Trunc(f) == f { // integral value
            if f >= float64(minInt) && f <= float64(maxInt) {
                return int(f)
            }
            if f >= float64(math.MinInt64) && f <= float64(math.MaxInt64) {
                return int64(f)
            }
        }
        return f
    case int:
        return val
    case int8:
        return int(val)
    case int16:
        return int(val)
    case int32:
        return int(val)
    case int64:
        if val >= minInt && val <= maxInt {
            return int(val)
        }
        return val
    case uint8:
        return int(val)
    case uint16:
        return int(val)
    case uint32:
        if uint64(val) <= uint64(maxInt) { // fits in signed int
            return int(val)
        }
        return int64(val)
    case uint:
        if uint64(val) <= uint64(maxInt) { // fits in signed int
            return int(val)
        }
        return int64(val)
    case uint64:
        u := val
        if u <= uint64(maxInt) {
            return int(u)
        }
        if u <= uint64(math.MaxInt64) {
            return int64(u)
        }
        return u // preserve large uint64 that can't fit in int64
    case string:
        if i, err := strconv.ParseInt(val, 10, 64); err == nil {
            if i >= minInt && i <= maxInt {
                return int(i)
            }
            return i
        }
        return v
    default:
        return v
    }
}

func platformIntLimits() (min int64, max int64) {
    if strconv.IntSize == 32 {
        return int64(-1 << 31), int64((1<<31)-1)
    }
    return math.MinInt64, math.MaxInt64
}
