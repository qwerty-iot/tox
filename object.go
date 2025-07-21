package tox

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/goccy/go-json"
	"math"
	"reflect"
	"strings"
	"time"
)

type Object map[string]any

func NewObject(mi any) Object {
	if mi == nil {
		mi = map[string]any{}
	}
	switch mt := mi.(type) {
	case map[string]any:
		o := Object(mt)
		o = o.Clone()
		o.ConvertStructs()
		return o
	case Object:
		no := mt.Clone()
		no.ConvertStructs()
		return no
	case []byte:
		var obj Object
		_ = json.Unmarshal(mt, &obj)
		return obj
	case string:
		var obj Object
		_ = json.Unmarshal([]byte(mt), &obj)
		return obj
	default:
		return structToObject(mi)
	}
}

func (o Object) Clone() Object {
	if o == nil {
		return nil
	}
	no, _ := Deepcopy(o)
	return no
}

func (o Object) ToStruct(target any) {
	b, _ := json.Marshal(o)
	_ = json.Unmarshal(b, target)
}

func (o Object) Equals(other Object) bool {
	if (o == nil && other != nil) || (o != nil && other == nil) {
		return false
	} else if o == nil {
		return true
	}
	if fmt.Sprintf("%v", o) == fmt.Sprintf("%v", other) {
		return true
	}
	return false
}

func countFields(x any) int {
	count := 0
	switch tx := x.(type) {
	case []any:
		for _, v := range tx {
			count += countFields(v)
		}
	case map[string]any:
		for _, v := range tx {
			count += countFields(v)
		}
	default:
		count++
	}
	return count
}

func (o Object) FieldCount() int {
	if o == nil {
		return 0
	}
	return countFields(o)
}

func keyIndex(key string) (string, int) {
	lb := strings.Index(key, "[")
	rb := strings.Index(key, "]")
	var sk string
	var si int
	if lb > 0 {
		sk = key[:lb]
		if rb == -1 {
			si = ToInt(key[lb+1:])
		} else {
			si = ToInt(key[lb+1 : rb])
		}
		return sk, si
	} else {
		return key, -1
	}
}

func (o Object) Move(from string, to string) {
	if o == nil {
		return
	}
	if from == to {
		return
	}
	f := o.Get(from)
	if f != nil {
		o.Set(to, f)
		o.Delete(from)
	}
}

func (o Object) Delete(key string) {
	if o == nil {
		return
	}
	idx := strings.LastIndex(key, ".")
	if idx == -1 {
		delete(o, key)
	} else {
		parent := o.GetObject(key[:idx])
		if parent != nil {
			delete(parent, key[idx+1:])
		}
	}
}

func (o Object) DeletePrefix(prefix string) {
	if o == nil {
		return
	}
	idx := strings.LastIndex(prefix, ".")
	if idx == -1 {
		for k := range o {
			if strings.HasPrefix(k, prefix) {
				delete(o, k)
			}
		}
	} else {
		parent := o.GetObject(prefix[:idx])
		if parent != nil {
			for k := range parent {
				if strings.HasPrefix(k, prefix) {
					delete(o, k)
				}
			}
		}
	}
}

func (o Object) Exists(key string) bool {
	if o == nil {
		return false
	}
	if o.Get(key) != nil {
		return true
	}
	return false
}

func (o Object) Get(key string) any {
	if o == nil {
		return nil
	}
	parts := strings.Split(key, ".")
	if len(parts) == 1 {
		k, i := keyIndex(key)
		if field, found := o[k]; found {
			if i == -1 {
				return field
			} else {
				if fa, ok := field.([]any); ok {
					if i < 0 || i >= len(fa) {
						return nil
					}
					return fa[i]
				} else {
					return nil
				}
			}
		} else {
			return nil
		}
	} else {
		mi := o
		for partIdx, part := range parts {
			k, i := keyIndex(part)
			if p, found := mi[k]; found {
				if i == -1 {
					if pi, ok := p.(map[string]any); ok && partIdx != len(parts)-1 {
						mi = pi
					} else if pi, ok := p.(Object); ok && partIdx != len(parts)-1 {
						mi = pi
					} else if partIdx == len(parts)-1 {
						return p
					} else {
						return nil
					}
				} else {
					switch pa := p.(type) {
					case []int:
						if i < 0 || i >= len(pa) {
							return nil
						}
						return pa[i]
					case []int64:
						if i < 0 || i >= len(pa) {
							return nil
						}
						return pa[i]
					case []float32:
						if i < 0 || i >= len(pa) {
							return nil
						}
						return pa[i]
					case []float64:
						if i < 0 || i >= len(pa) {
							return nil
						}
						return pa[i]
					case []bool:
						if i < 0 || i >= len(pa) {
							return nil
						}
						return pa[i]
					case []string:
						if i < 0 || i >= len(pa) {
							return nil
						}
						return pa[i]
					case []map[string]any:
						if i < 0 || i >= len(pa) {
							return nil
						}
						if partIdx == len(parts)-1 {
							return pa[i]
						}
						mi = pa[i]
					case []Object:
						if i < 0 || i >= len(pa) {
							return nil
						}
						if partIdx == len(parts)-1 {
							return pa[i]
						}
						mi = pa[i]
					case []any:
						if i < 0 || i >= len(pa) {
							return nil
						}
						if partIdx == len(parts)-1 {
							return pa[i]
						}
						// Check if the array element is a map or Object and continue processing
						if mapVal, ok := pa[i].(map[string]any); ok {
							mi = mapVal
						} else if objVal, ok := pa[i].(Object); ok {
							mi = objVal
						} else {
							return pa[i]
						}
					default:
						return nil
					}
				}
			} else {
				return nil
			}
		}
	}
	return nil
}

func (o Object) GetObjectArray(key string) []Object {
	if o == nil {
		return nil
	}
	if field := o.Get(key); field != nil {
		fieldVal := reflect.ValueOf(field)
		if fieldVal.Kind() == reflect.Array || fieldVal.Kind() == reflect.Slice {
			var ret = make([]Object, fieldVal.Len())
			for i := 0; i < fieldVal.Len(); i++ {
				ret[i] = NewObject(fieldVal.Index(i).Interface())
			}
			return ret
		}
		return nil
	} else {
		return nil
	}
}

func (o Object) GetObject(key string) Object {
	if o == nil {
		return nil
	}
	if field := o.Get(key); field != nil {
		switch typed := field.(type) {
		case map[string]any:
			return typed
		case Object:
			return typed
		default:
			return nil
		}
	} else {
		return nil
	}
}

func (o Object) GetString(key string, def string) string {
	if o == nil {
		return def
	}
	if field := o.Get(key); field != nil {
		return ToString(field)
	} else {
		return def
	}
}

func (o Object) GetStringArray(key string, def []string) []string {
	if o == nil {
		return def
	}
	if field := o.Get(key); field != nil {
		return ToStringArray(field)
	} else {
		return def
	}
}

func (o Object) GetInt(key string, def int) int {
	if o == nil {
		return def
	}
	if field := o.Get(key); field != nil {
		return ToInt(field)
	} else {
		return def
	}
}

func (o Object) ModifyInt(key string, delta int) int {
	if o == nil {
		return delta
	}
	if field := o.Get(key); field != nil {
		v := ToInt(field) + delta
		o.Set(key, v)
		return v
	} else {
		o.Set(key, delta)
		return delta
	}
}

func (o Object) GetIntPtr(key string, def int) *int {
	if o == nil {
		return ToIntPtr(def)
	}
	if field := o.Get(key); field != nil {
		return ToIntPtr(field)
	} else {
		if def == math.MinInt {
			return nil
		}
		return ToIntPtr(def)
	}
}

func (o Object) GetFloat64(key string, def float64) float64 {
	if o == nil {
		return def
	}
	if field := o.Get(key); field != nil {
		return ToFloat64(field)
	} else {
		return def
	}
}

func (o Object) GetFloat64Ptr(key string, def float64) *float64 {
	if o == nil {
		return ToFloat64Ptr(def)
	}
	if field := o.Get(key); field != nil {
		return ToFloat64Ptr(field)
	} else {
		if def == math.NaN() {
			return nil
		}
		return ToFloat64Ptr(def)
	}
}

func (o Object) GetBool(key string, def bool) bool {
	if o == nil {
		return def
	}
	if field := o.Get(key); field != nil {
		return ToBool(field)
	} else {
		return def
	}
}

func (o Object) GetBoolPtr(key string, def bool) *bool {
	if o == nil {
		return ToBoolPtr(def)
	}
	if field := o.Get(key); field != nil {
		return ToBoolPtr(field)
	} else {
		return ToBoolPtr(def)
	}
}

func (o Object) GetTime(key string, def time.Time) time.Time {
	if o == nil {
		return def
	}
	if field := o.Get(key); field != nil {
		return ToTime(field)
	} else {
		return def
	}
}

func (o Object) GetTimePtr(key string, def time.Time) *time.Time {
	if o == nil {
		return ToTimePtr(def)
	}
	if field := o.Get(key); field != nil {
		return ToTimePtr(field)
	} else {
		if def.IsZero() {
			return nil
		}
		return ToTimePtr(def)
	}
}

func (o Object) GetBytes(key string, def []byte) []byte {
	if o == nil {
		return def
	}
	if field := o.Get(key); field != nil {
		switch typed := field.(type) {
		case []byte:
			return typed
		case string:
			return []byte(typed)
		default:
			return []byte(ToJson(typed))
		}
	} else {
		return def
	}
}

func (o Object) GetInto(key string, ret any) any {
	if o == nil {
		return ret
	}
	if field := o.Get(key); field != nil {
		b, _ := json.Marshal(field)
		_ = json.Unmarshal(b, ret)
		return ret
	} else {
		return ret
	}
}

func (o Object) Unmarshal(field string, raw any) {
	if o == nil {
		return
	}
	var parsed bool
	switch v := raw.(type) {
	case []byte:
		if len(v) == 0 {
			return
		} else if v[0] == '{' {
			var data Object
			err := json.Unmarshal(v, &data)
			if err == nil {
				o[field] = data
				parsed = true
			}
		} else if v[0] == '[' {
			var data []Object
			err := json.Unmarshal(v, &data)
			if err == nil {
				o[field] = data
				parsed = true
			}
		}
		if !parsed {
			if isUnicode(v) {
				o[field] = string(v)
			} else {
				o[field] = raw
			}
		}
	case string:
		if len(v) == 0 {
			return
		} else if v[0] == '{' {
			var data Object
			err := json.Unmarshal([]byte(v), &data)
			if err == nil {
				o[field] = data
				parsed = true
			}
		} else if v[0] == '[' {
			var data []Object
			err := json.Unmarshal([]byte(v), &data)
			if err == nil {
				o[field] = data
				parsed = true
			}
		}
		if !parsed {
			o[field] = v
		}
	case map[string]any:
		o[field] = v
	case []map[string]any:
		o[field] = v
	default:
		o[field] = structToObject(v)
	}

}

const (
	BlankString = "(blank)"
	Null        = "(null)"
)

func (o Object) SetNullIfNotExist(key string) {
	if o == nil {
		return
	}
	if o.Get(key) == nil {
		o.Set(key, Null)
	}
}

func (o Object) SetIfNotExist(key string, value any) {
	if o == nil {
		return
	}
	if o.Get(key) == nil {
		o.Set(key, value)
	}
}

func (o Object) Set(key string, value any) {
	if value == nil {
		return
	} else if vs, ok := value.(string); ok && len(vs) == 0 {
		return
	} else if vm, ok := value.(map[string]any); ok && len(vm) == 0 {
		return
	}
	if value == BlankString {
		value = ""
	} else if value == Null {
		value = nil
	}
	parts := strings.Split(key, ".")
	if len(parts) == 1 {
		o[key] = value
	} else {
		mi := o
		for partIdx, part := range parts {
			if p, found := mi[part]; found {
				if pi, ok := p.(map[string]any); ok && partIdx != len(parts)-1 {
					mi = pi
				} else if pi, ok := p.(Object); ok && partIdx != len(parts)-1 {
					mi = pi
				} else {
					mi[part] = value
				}
			} else {
				if partIdx != len(parts)-1 {
					pi := Object{}
					mi[part] = pi
					mi = pi
				} else {
					mi[part] = value
				}
			}
		}
	}
	return
}

func (o Object) Merge(other Object) {
	_ = mergo.Merge(&o, other, mergo.WithOverride)
}
func (o Object) MergeMissing(other Object) {
	_ = mergo.Merge(&o, other)
}

func (o Object) Flatten(delim string) Object {
	output := make(map[string]any)

	var flatten func(map[string]any, string)
	flatten = func(m map[string]any, parentKey string) {
		for k, v := range m {
			key := parentKey + k
			if parentKey != "" {
				key = parentKey + delim + k
			}

			switch value := v.(type) {
			case map[string]any:
				flatten(value, key)
			case Object:
				flatten(value, key)
			default:
				fieldVal := reflect.ValueOf(value)
				if fieldVal.Kind() == reflect.Array || fieldVal.Kind() == reflect.Slice {
					for i := 0; i < fieldVal.Len(); i++ {
						vv := fieldVal.Index(i).Interface()
						switch vvv := vv.(type) {
						case map[string]any:
							flatten(vvv, key+fmt.Sprintf("[%d]", i))
						default:
							output[key+fmt.Sprintf("[%d]", i)] = vvv
						}
					}
				} else {
					output[key] = value
				}
			}
		}
	}

	flatten(o, "")

	return output
}

type FieldDiff struct {
	Old any `json:"old" bson:"old"`
	New any `json:"new" bson:"new"`
}
type ObjectDiff struct {
	Same     bool                 `json:"same"               bson:"same"`
	Added    Object               `json:"added,omitempty"    bson:"added,omitempty"`
	Modified map[string]FieldDiff `json:"modified,omitempty" bson:"modified,omitempty"`
	Deleted  Object               `json:"deleted,omitempty"  bson:"deleted,omitempty"`
}

func (o Object) JsonString(pretty bool) string {
	if o == nil {
		return ""
	}
	if pretty {
		b, _ := json.MarshalIndent(o, "", "  ")
		return string(b)
	} else {
		b, _ := json.Marshal(o)
		return string(b)
	}
}

func (o Object) JsonBytes(pretty bool) []byte {
	if o == nil {
		return nil
	}
	if pretty {
		b, _ := json.MarshalIndent(o, "", "  ")
		return b
	} else {
		b, _ := json.Marshal(o)
		return b
	}
}

func (o Object) Diff(other Object) ObjectDiff {
	if o == nil && other == nil {
		return ObjectDiff{Same: true}
	} else if o == nil || other == nil {
		return ObjectDiff{Same: false}
	}
	origf := o.Flatten("/")
	otherf := other.Flatten("/")

	diff := diffMaps(origf, otherf)
	if diff.Added == nil && diff.Modified == nil && diff.Deleted == nil {
		diff.Same = true
	} else {
		diff.Same = false
	}

	return diff
}

func (o Object) RemoveNaN() {
	if o == nil {
		return
	}
	toBeDeleted := []string{}
	o.ConvertStructs()
	removeNaN(o, "", &toBeDeleted)
	if len(toBeDeleted) > 0 {
		for _, key := range toBeDeleted {
			o.Delete(key)
		}
	}
}

func (o Object) ConvertStructs() {
	if o == nil {
		return
	}
	b, _ := json.Marshal(o)
	_ = json.Unmarshal(b, &o)
}
