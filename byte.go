package tox

// ToBool converts bool, string, or byte arrays, if the conversion fails, it returns false.
func ToByteArray(v interface{}) []byte {
	switch v := v.(type) {
	case bool:
		if v {
			return []byte{0x01}
		} else {
			return []byte{0x00}
		}
	case nil:
		return nil
	case string:
		return []byte(v)
	case []byte:
		return v
	default:
		return nil
	}
}
