package tox

import "strconv"

func ToByte(v interface{}) byte {
	switch v := v.(type) {
	case byte:
		return v
	case int:
		return byte(v)
	case int8:
		return byte(v)
	case int16:
		return byte(v)
	case int32:
		return byte(v)
	case int64:
		return byte(v)
	case uint:
		return byte(v)
	case uint16:
		return byte(v)
	case uint32:
		return byte(v)
	case uint64:
		return byte(v)
	case float32:
		return byte(v)
	case float64:
		return byte(v)
	case string:
		i, _ := strconv.Atoi(v)
		return byte(i)
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// ToByteArray converts bool, string, or byte arrays, if the conversion fails, it returns false.
func ToByteArray(v interface{}) []byte {
	switch v := v.(type) {
	case bool:
		if v {
			return []byte{0x01}
		}
		return []byte{0x00}
	case nil:
		return nil
	case string:
		return []byte(v)
	case []byte:
		return v
	case []any:
		rslt := make([]byte, len(v))
		for i, v := range v {
			rslt[i] = ToByte(v)
		}
		return rslt
	default:
		return nil
	}
}

func TruncateByteArray(b []byte, length int) []byte {
	if len(b) > length {
		return b[:length]
	}
	return b
}
