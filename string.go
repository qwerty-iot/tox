package tox

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/goccy/go-json"
)

func toStringMap(in map[any]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		sk := ToString(k)
		switch vv := v.(type) {
		case map[any]any:
			out[sk] = toStringMap(vv)
		default:
			out[sk] = v
		}
	}
	return out
}

// ToString converts any data type to a string, it uses fmt.Sprintf() to convert unknown types.
func ToString(v interface{}) string {
	return ToStringOpts(v, nil)
}
func ToStringOpts(v interface{}, options *Options) string {
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
	case float64:
		if options != nil && options.FloatToInt {
			if v == math.Floor(v) {
				return fmt.Sprintf("%.0f", v)
			}
		}
		if options != nil && options.FloatPrecision > 0 {
			return fmt.Sprintf("%.*f", options.FloatPrecision, v)
		} else {
			return fmt.Sprintf("%v", v)
		}
	case float32:
		return ToStringOpts(float64(v), options)
	case int, int64, uint, uint64, int8, int16, uint8, uint16:
		return fmt.Sprintf("%v", v)
	case time.Duration:
		return ToPrettyDuration(v, FormatShort)
	case time.Time:
		return v.Format(time.RFC3339Nano)
	case []byte:
		if utf8.Valid(v) {
			return string(v)
		} else {
			return fmt.Sprintf("%v", v)
		}
	case map[string]any:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		} else {
			return string(b)
		}
	case map[any]any:
		b, err := json.Marshal(v)
		if err != nil {
			var typeErr *json.UnsupportedTypeError
			if errors.As(err, &typeErr) {
				fm := toStringMap(v)
				return ToJson(fm)
			}
			return fmt.Sprintf("%v", v)
		} else {
			return string(b)
		}
	default:

		if reflect.ValueOf(v).Kind() == reflect.Ptr {
			if reflect.ValueOf(v).IsNil() {
				return ""
			} else {
				return ToString(reflect.ValueOf(v).Elem().Interface())
			}
		}

		switch reflect.TypeOf(v).Name() {
		case "DateTime":
			if reflect.TypeOf(v).Kind() == reflect.Int64 {
				return time.Unix(reflect.ValueOf(v).Int()/1000, 0).Format(time.RFC3339Nano)
			}
		}
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

func ToPrettyJson(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
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

type StringFormat int

const (
	FormatShort  StringFormat = 1
	FormatMedium StringFormat = 2
	FormatLong   StringFormat = 3
)

func ToPrettyDuration(v time.Duration, format StringFormat) string {
	var r []string
	isNegative := false
	if v < 0 {
		v = -v
		isNegative = true
	}
	days := v / (24 * time.Hour)
	if days > 0 {
		switch format {
		case FormatMedium:
			r = append(r, fmt.Sprintf("%dd", days))
		case FormatLong:
			r = append(r, fmt.Sprintf("%d days", days))
		}
	}
	v -= days * 24 * time.Hour

	hours := v / time.Hour
	if hours > 0 {
		switch format {
		case FormatMedium:
			r = append(r, fmt.Sprintf("%dh", hours))
		case FormatLong:
			r = append(r, fmt.Sprintf("%d hours", hours))
		}
	}
	v -= hours * time.Hour

	minutes := v / time.Minute
	if minutes > 0 {
		switch format {
		case FormatMedium:
			r = append(r, fmt.Sprintf("%dm", minutes))
		case FormatLong:
			r = append(r, fmt.Sprintf("%d minutes", minutes))
		}
	}
	v -= minutes * time.Minute

	seconds := v / time.Second
	if seconds > 0 {
		switch format {
		case FormatMedium:
			r = append(r, fmt.Sprintf("%ds", seconds))
		case FormatLong:
			r = append(r, fmt.Sprintf("%d seconds", seconds))
		}
	}

	var ret string
	if format == FormatShort {
		ret = fmt.Sprintf("%d.%d:%d:%d", days, hours, minutes, seconds)
	} else {
		ret = strings.Join(r, ", ")
	}
	if isNegative {
		ret = "-" + ret
	}
	return ret
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

func CapitalizeString(s string) string {
	// Split the input string into words
	words := strings.Fields(s)

	// Initialize an empty slice to store the capitalized words
	capitalizedWords := make([]string, len(words))

	// Capitalize the first letter of each word and make the rest lowercase
	for i, word := range words {
		// Check if the word is empty
		if len(word) == 0 {
			continue
		}

		// Convert the first letter to uppercase and the rest to lowercase
		capitalizedWord := string(unicode.ToUpper(rune(word[0]))) + strings.ToLower(word[1:])
		capitalizedWords[i] = capitalizedWord
	}

	// Join the capitalized words back into a single string
	result := strings.Join(capitalizedWords, " ")

	return result
}

func StringInArray(s string, arr []string, ignoreCase bool) bool {
	if arr == nil {
		return false
	}
	for _, v := range arr {
		if ignoreCase {
			if strings.EqualFold(v, s) {
				return true
			}
		} else {
			if v == s {
				return true
			}
		}
	}
	return false
}
