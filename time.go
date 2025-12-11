package tox

import (
	"reflect"
	"strconv"
	"strings"
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

func ParseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, strconv.ErrSyntax
	}
	// support suffixes: ms, s, m, h, d
	// try ms explicitly
	if strings.HasSuffix(s, "ms") {
		v := strings.TrimSuffix(s, "ms")
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Millisecond * time.Duration(n), nil
	}
	// single-char suffixes
	unit := s[len(s)-1]
	num := s[:len(s)-1]
	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 0, err
	}
	switch unit {
	case 's':
		return time.Second * time.Duration(n), nil
	case 'm':
		return time.Minute * time.Duration(n), nil
	case 'h':
		return time.Hour * time.Duration(n), nil
	case 'd':
		return time.Hour * time.Duration(n*24), nil
	default:
		// if no recognized suffix, try to parse as plain milliseconds
		n2, err2 := strconv.ParseInt(s, 10, 64)
		if err2 != nil {
			return 0, strconv.ErrSyntax
		}
		return time.Second * time.Duration(n2), nil
	}
}
