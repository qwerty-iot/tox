package tox

// ToMapStringString converts a generic map[string]interface{} to a map[string]string.
func ToMapStringString(v interface{}) map[string]string {
	switch v := v.(type) {
	case map[string]interface{}:
		ret := map[string]string{}
		for key, val := range v {
			ret[key] = ToString(val)
		}
		return ret
	default:
		return nil
	}
}
