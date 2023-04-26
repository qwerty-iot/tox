package tox

import "encoding/json"

// ToMapStringString converts a generic map[string]interface{} to a map[string]string.
func ToMapStringString(v interface{}) map[string]string {
	switch v := v.(type) {
	case map[string]string:
		return v
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

func ToMapStringInterface(v interface{}) map[string]interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		return v
	case map[string]string:
		ret := map[string]interface{}{}
		for key, val := range v {
			ret[key] = val
		}
		return ret
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		ret := map[string]interface{}{}
		err = json.Unmarshal(b, &ret)
		if err != nil {
			return nil
		}
		return ret
	}
}

func FlattenMap(m map[string]any, delim string) map[string]any {
	output := make(map[string]any)

	hasSubmaps := false
	for _, value := range m {
		if _, ok := value.(map[string]any); ok {
			hasSubmaps = true
		}
	}
	if !hasSubmaps {
		return m
	}
	for key, value := range m {
		if submap, ok := value.(map[string]any); ok {
			flattenedSubmap := FlattenMap(submap, delim)
			for subkey, subvalue := range flattenedSubmap {
				output[key+delim+subkey] = subvalue
			}
		} else {
			output[key] = value
		}
	}
	return output
}
