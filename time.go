package tox

import "time"

func ToTime(v interface{}) time.Time {
	switch v := v.(type) {
	case int:
		return time.Unix(int64(v), 0)
	case int64:
		return time.Unix(v, 0)
	case time.Time:
		return v
	case string:
		t, err := time.Parse(time.RFC3339Nano, v)
		if err != nil {
			return time.Time{}
		}
		return t
	default:
		return time.Time{}
	}
}
