package tox

import (
	"fmt"
)

func ToString(v interface{}) string {
	switch v := v.(type) {
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	case nil:
		return ""
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

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
