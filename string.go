package tox

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode/utf8"
)

// ToString converts any data type to a string, it uses fmt.Sprintf() to convert unknown types.
func ToString(v interface{}) string {
	switch v := v.(type) {
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return ""
	case string:
		return v
	case int, int64, uint, uint64, float64, float32, int8, int16, uint8, uint16:
		return fmt.Sprintf("%v", v)
	case []byte:
		if utf8.Valid(v) {
			return string(v)
		} else {
			return fmt.Sprintf("%v", v)
		}
	case map[string]interface{}:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		} else {
			return string(b)
		}
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		} else {
			return string(b)
		}
	}
}

func ToJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err == nil {
		return string(b)
	} else {
		return fmt.Sprintf("%v", v)
	}
}

// ToStringArray can convert a single string to an array, useful if interface could be a string or array of strings.
func ToStringArray(v interface{}) []string {
	switch v := v.(type) {
	case nil:
		return nil
	case string:
		return []string{v}
	case []string:
		return v
	case []any:
		var ret = make([]string, len(v))
		for ii, vv := range v {
			ret[ii] = ToString(vv)
		}
		return ret
	default:
		aVal := reflect.ValueOf(v)
		if aVal.Kind() == reflect.Array || aVal.Kind() == reflect.Slice {
			var ret = make([]string, aVal.Len())
			for i := 0; i < aVal.Len(); i++ {
				ret[i] = ToString(aVal.Index(i).Interface())
			}
			return ret
		}
		return nil
	}
}

func ToStringPtr(v interface{}) *string {
	if v == nil {
		return nil
	}
	ret := ToString(v)
	return &ret
}

func ToStringPtrOpts(v interface{}, options *Options) *string {
	if v == nil {
		return nil
	}
	ret := ToString(v)
	if options != nil && options.EmptyStringAsNull && ret == "" {
		return nil
	}
	return &ret
}

func TruncateString(s string, length int) string {
	if len(s) > length {
		return s[:length]
	}
	return s
}
