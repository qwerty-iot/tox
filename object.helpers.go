package tox

import (
	"github.com/goccy/go-json"
	"math"
	"reflect"
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

func structToObject(input any) Object {
	b, err := json.Marshal(input)
	if err == nil {
		var obj Object
		_ = json.Unmarshal(b, &obj)
		return obj
	}
	return nil
}

/*
future implementation
func structToObject(input any) any {
	structValue := reflect.ValueOf(input)
	if structValue.Kind() == reflect.Struct {
		result := Object{}
		structType := reflect.TypeOf(input)
		for i := 0; i < structValue.NumField(); i++ {
			field := structValue.Field(i)
			fieldName := structType.Field(i).Name
			if jt := structType.Field(i).Tag.Get("json"); len(jt) != 0 {
				ss := strings.Split(jt, ",")
				fieldName = ss[0]
			}
			if !unicode.IsLower(rune(fieldName[0])) {
				result[fieldName] = structToObject(field.Interface())
			}
		}
		return result
	} else if structValue.Kind() == reflect.Ptr && structValue.Elem().Kind() == reflect.Struct {
		return structToObject(structValue.Elem().Interface())
	}
	return input
}
*/

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
