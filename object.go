package tox

import (
	"encoding/json"
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/imdario/mergo"
)

type Object map[string]any

func NewObject(mi any) Object {
	if mi == nil {
		mi = map[string]any{}
	}
	switch mt := mi.(type) {
	case map[string]any:
		return Object(mt)
	case []byte:
		var obj Object
		_ = json.Unmarshal(mt, &obj)
		return obj
	case string:
		var obj Object
		_ = json.Unmarshal([]byte(mt), &obj)
		return obj
	default:
		b, err := json.Marshal(mt)
		if err == nil {
			var obj Object
			_ = json.Unmarshal(b, &obj)
			return obj
		}
		return nil
	}
	return nil
}

func (o Object) Clone() Object {
	no, _ := Deepcopy(o)
	return no.(Object)
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

func (o Object) Delete(key string) {
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

func (o Object) Get(key string) any {
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
					} else {
						return p
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
						return pa[i]
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
	if field := o.Get(key); field != nil {
		if arr, ok := field.([]any); ok {
			var ret []Object
			for _, item := range arr {
				if mi, ok := item.(map[string]any); ok {
					ret = append(ret, mi)
				}
			}
			return ret
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (o Object) GetObject(key string) Object {
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
	if field := o.Get(key); field != nil {
		return ToString(field)
	} else {
		return def
	}
}

func (o Object) GetStringArray(key string, def []string) []string {
	if field := o.Get(key); field != nil {
		return ToStringArray(field)
	} else {
		return def
	}
}

func (o Object) GetInt(key string, def int) int {
	if field := o.Get(key); field != nil {
		return ToInt(field)
	} else {
		return def
	}
}

func (o Object) GetIntPtr(key string, def int) *int {
	if field := o.Get(key); field != nil {
		return ToIntPtr(field)
	} else {
		if def == math.MinInt {
			return nil
		}
		return ToIntPtr(def)
	}
}

func (o Object) GetBool(key string, def bool) bool {
	if field := o.Get(key); field != nil {
		return ToBool(field)
	} else {
		return def
	}
}

func (o Object) GetBoolPtr(key string, def bool) *bool {
	if field := o.Get(key); field != nil {
		return ToBoolPtr(field)
	} else {
		return ToBoolPtr(def)
	}
}

func (o Object) GetTime(key string, def time.Time) time.Time {
	if field := o.Get(key); field != nil {
		return ToTime(field)
	} else {
		return def
	}
}

func (o Object) GetTimePtr(key string, def time.Time) *time.Time {
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

func (o Object) Unmarshal(field string, raw any) {
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
			if isASCII(v) {
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
		o[field] = v
	}

}

func isASCII(s []byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		} else if s[i] < 0x20 && s[i] != 0x09 && s[i] != 0x0A && s[i] != 0x0D {
			return false
		}
	}
	return true
}

func (o Object) Set(key string, value any) {
	if value == nil {
		return
	} else if vs, ok := value.(string); ok && len(vs) == 0 {
		return
	} else if vm, ok := value.(map[string]any); ok && len(vm) == 0 {
		return
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
	_ = mergo.Merge(&o, other)
}
