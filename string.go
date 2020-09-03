package tox

import (
	"fmt"
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
	default:
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
	default:
		return nil
	}
}
