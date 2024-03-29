package tox

import (
	"reflect"
	"time"
)

// ToTime converts data types to time.Time structures.  'int' or 'int64' are treated as unix time, strings are treated
// as RFC3330Nano timestamps.  If the conversion fails, an empty time.Time{} is returned.
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
		if reflect.TypeOf(v).Kind() == reflect.Int64 && reflect.TypeOf(v).Name() == "DateTime" {
			// special case for mongodb primitive.DateTime
			return time.Unix(reflect.ValueOf(v).Int()/1000, 0)
		}
		return time.Time{}
	}
}

func ToTimePtr(v interface{}) *time.Time {
	if v == nil {
		return nil
	}
	ret := ToTime(v)
	return &ret
}
