package tox

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
	"unicode"
)

func isUnicode(s []byte) bool {
	for _, r := range string(s) {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func isASCII(s []byte) bool {
	for i := 0; i < len(s); i++ {
		// Check if character is above ASCII range
		if s[i] > 165 {
			return false
		}
		// Allow printable characters (0x20-0x7E)
		if s[i] >= 0x20 && s[i] <= 0x7E {
			continue
		}
		// Allow specific control characters: tab, newline, carriage return, space
		if s[i] != 0x09 && s[i] != 0x0A && s[i] != 0x0D && s[i] != 0x20 {
			return false
		}
	}
	return true
}

func structToObject(input any) Object {
	res := structToAnything(input)
	if obj, ok := res.(Object); ok {
		return obj
	}
	if m, ok := res.(map[string]any); ok {
		return Object(m)
	}
	return nil
}

func structToAnything(input any) any {
	if input == nil {
		return nil
	}

	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	// Handle pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		t = v.Type()
	}

	// Special case: time.Time should be returned as is
	if _, ok := v.Interface().(time.Time); ok {
		return v.Interface()
	}

	// Handle types with Hex() method that look like ObjectID
	if strings.Contains(t.String(), "ObjectID") {
		if m, ok := v.Interface().(interface{ Hex() string }); ok {
			return m.Hex()
		}
	}

	// Handle pointers to types with Hex() method that look like ObjectID
	if t.Kind() == reflect.Ptr && strings.Contains(t.Elem().String(), "ObjectID") {
		if !v.IsNil() {
			if m, ok := v.Interface().(interface{ Hex() string }); ok {
				return m.Hex()
			}
		}
	}

	// Special case: []byte should be returned as is
	if b, ok := v.Interface().([]byte); ok {
		return b
	}

	switch v.Kind() {
	case reflect.Struct:
		result := Object{}
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			structField := t.Field(i)

			// Skip unexported fields
			if structField.PkgPath != "" {
				continue
			}

			fieldName := structField.Name
			omitempty := false
			if jt := structField.Tag.Get("json"); len(jt) != 0 {
				if jt == "-" {
					continue
				}
				ss := strings.Split(jt, ",")
				if ss[0] != "" {
					fieldName = ss[0]
				}
				for _, opt := range ss[1:] {
					if opt == "omitempty" {
						omitempty = true
					}
				}
			}

			if omitempty && field.IsZero() {
				continue
			}

			result[fieldName] = structToAnything(field.Interface())
		}
		return result

	case reflect.Slice, reflect.Array:
		result := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = structToAnything(v.Index(i).Interface())
		}
		return result

	case reflect.Map:
		result := make(map[string]any)
		for _, key := range v.MapKeys() {
			result[fmt.Sprintf("%v", key.Interface())] = structToAnything(v.MapIndex(key).Interface())
		}
		return result

	default:
		return v.Interface()
	}
}

func removeNaN(a any, parent string, toBeDeleted *[]string) {
	if len(parent) != 0 {
		parent += "."
	}
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			originalValue := v.MapIndex(key)
			if !isNan(originalValue.Interface()) {
				removeNaN(originalValue.Interface(), parent+key.String(), toBeDeleted)
			} else {
				*toBeDeleted = append(*toBeDeleted, parent+key.String())
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			originalValue := v.Index(i)
			if isNan(originalValue.Interface()) {
				v.Index(i).Set(reflect.ValueOf(float64(0.0)))
			}
		}

		//case reflect.Struct:
		//	spew.Dump("struct", v.Interface())
		//	panic("no structs allowed in removeNaN")
		//case reflect.Ptr:
		//	panic("no struct pointers allowed in removeNaN")
	default:
		// do nothing
	}
}

func isNan(a any) bool {
	v := reflect.ValueOf(a)
	return v.Kind() == reflect.Float64 && math.IsNaN(v.Float())
}

func diffMaps(oldMap, newMap map[string]any) ObjectDiff {
	result := ObjectDiff{
		Added:    Object{},
		Modified: make(map[string]FieldDiff),
		Deleted:  Object{},
	}

	// Check for added and modified key-value pairs
	for key, newValue := range newMap {
		if oldValue, exists := oldMap[key]; exists {
			if !reflect.DeepEqual(oldValue, newValue) {
				result.Modified[key] = FieldDiff{Old: oldValue, New: newValue}
			}
		} else {
			result.Added[key] = newValue
		}
	}

	// Check for deleted key-value pairs
	for key, oldValue := range oldMap {
		if _, exists := newMap[key]; !exists {
			result.Deleted[key] = oldValue
		}
	}
	if len(result.Added) == 0 {
		result.Added = nil
	}
	if len(result.Deleted) == 0 {
		result.Deleted = nil
	}
	if len(result.Modified) == 0 {
		result.Modified = nil
	}

	return result
}
