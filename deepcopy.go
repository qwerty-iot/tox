/*

MIT License

Copyright (c) 2020 Sabrina Williams

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package tox

import (
	"errors"
	"fmt"
	. "reflect"
)

type copier func(any, map[uintptr]any) (any, error)

var copiers map[Kind]copier

func init() {
	copiers = map[Kind]copier{
		Bool:       _primitive,
		Int:        _primitive,
		Int8:       _primitive,
		Int16:      _primitive,
		Int32:      _primitive,
		Int64:      _primitive,
		Uint:       _primitive,
		Uint8:      _primitive,
		Uint16:     _primitive,
		Uint32:     _primitive,
		Uint64:     _primitive,
		Uintptr:    _primitive,
		Float32:    _primitive,
		Float64:    _primitive,
		Complex64:  _primitive,
		Complex128: _primitive,
		Array:      _array,
		Map:        _map,
		Ptr:        _pointer,
		Slice:      _slice,
		String:     _primitive,
		Struct:     _struct,
	}
}

// Primitive makes a copy of a primitive type...which just means it returns the input value.
// This is wholly uninteresting, but I included it for consistency's sake.
func _primitive(x any, ptrs map[uintptr]any) (any, error) {
	kind := ValueOf(x).Kind()
	if kind == Array || kind == Chan || kind == Func || kind == Interface || kind == Map || kind == Ptr || kind == Slice || kind == Struct || kind == UnsafePointer {
		return nil, fmt.Errorf("unable to copy %v (a %v) as a primitive", x, kind)
	}
	return x, nil
}

// Deepcopy makes a deep copy of whatever gets passed in. It handles pretty much all known Go types
// (with the exception of channels, unsafe pointers, and functions). Note that this is a truly deep
// copy that will work it's way all the way to the leaves of the types--any pointer will be copied,
// any values in any slice or map will be deep copied, etc.
// Note: in order to avoid an infinite loop, we keep track of any pointers that we've run across.
// If we run into that pointer again, we don't make another deep copy of it; we just replace it with
// the copy we've already made. This also ensures that the cloned result is functionally equivalent
// to the original value.
func Deepcopy[T any](x T) (T, error) {
	ptrs := make(map[uintptr]any)
	ret, err := _anything(x, ptrs)
	if err != nil {
		return x, err
	}
	if rett, ok := ret.(T); ok {
		return rett, nil
	} else {
		return x, errors.New("unable to cast return value")
	}
}

func _anything(x any, ptrs map[uintptr]any) (any, error) {
	v := ValueOf(x)
	if !v.IsValid() {
		return x, nil
	}
	if c, ok := copiers[v.Kind()]; ok {
		return c(x, ptrs)
	}
	t := TypeOf(x)
	return nil, fmt.Errorf("unable to make a deep copy of %v (type: %v) - kind %v is not supported", x, t, v.Kind())
}

func _slice(x any, ptrs map[uintptr]any) (any, error) {
	v := ValueOf(x)
	if v.Kind() != Slice {
		return nil, fmt.Errorf("must pass a value with kind of Slice; got %v", v.Kind())
	}
	// Create a new slice and, for each item in the slice, make a deep copy of it.
	size := v.Len()
	t := TypeOf(x)
	dc := MakeSlice(t, size, size)
	for i := 0; i < size; i++ {
		item, err := _anything(v.Index(i).Interface(), ptrs)
		if err != nil {
			return nil, fmt.Errorf("failed to clone slice item at index %v: %v", i, err)
		}
		iv := ValueOf(item)
		if iv.IsValid() {
			dc.Index(i).Set(iv)
		}
	}
	return dc.Interface(), nil
}

func _map(x any, ptrs map[uintptr]any) (any, error) {
	v := ValueOf(x)
	if v.Kind() != Map {
		return nil, fmt.Errorf("must pass a value with kind of Map; got %v", v.Kind())
	}
	t := TypeOf(x)
	dc := MakeMapWithSize(t, v.Len())
	iter := v.MapRange()
	for iter.Next() {
		item, err := _anything(iter.Value().Interface(), ptrs)
		if err != nil {
			return nil, fmt.Errorf("failed to clone map item %v: %v", iter.Key().Interface(), err)
		}
		k, err := _anything(iter.Key().Interface(), ptrs)
		if err != nil {
			return nil, fmt.Errorf("failed to clone the map key %v: %v", k, err)
		}
		if item == nil {
			if mi, ok := dc.Interface().(map[string]any); ok {
				mi[ValueOf(k).String()] = nil
			}
		} else {
			dc.SetMapIndex(ValueOf(k), ValueOf(item))
		}
	}
	return dc.Interface(), nil
}

func _pointer(x any, ptrs map[uintptr]any) (any, error) {
	v := ValueOf(x)
	if v.Kind() != Ptr {
		return nil, fmt.Errorf("must pass a value with kind of Ptr; got %v", v.Kind())
	}

	if v.IsNil() {
		t := TypeOf(x)
		return Zero(t).Interface(), nil
	}

	addr := v.Pointer()
	if dc, ok := ptrs[addr]; ok {
		return dc, nil
	}
	t := TypeOf(x)
	dc := New(t.Elem())
	ptrs[addr] = dc.Interface()
	if !v.IsNil() {
		item, err := _anything(v.Elem().Interface(), ptrs)
		if err != nil {
			return nil, fmt.Errorf("failed to copy the value under the pointer %v: %v", v, err)
		}
		iv := ValueOf(item)
		if iv.IsValid() {
			dc.Elem().Set(ValueOf(item))
		}
	}
	return dc.Interface(), nil
}

func _struct(x any, ptrs map[uintptr]any) (any, error) {
	v := ValueOf(x)
	if v.Kind() != Struct {
		return nil, fmt.Errorf("must pass a value with kind of Struct; got %v", v.Kind())
	}
	t := TypeOf(x)
	dc := New(t)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		item, err := _anything(v.Field(i).Interface(), ptrs)
		if err != nil {
			return nil, fmt.Errorf("failed to copy the field %v in the struct %#v: %v", t.Field(i).Name, x, err)
		}
		dc.Elem().Field(i).Set(ValueOf(item))
	}
	return dc.Elem().Interface(), nil
}

func _array(x any, ptrs map[uintptr]any) (any, error) {
	v := ValueOf(x)
	if v.Kind() != Array {
		return nil, fmt.Errorf("must pass a value with kind of Array; got %v", v.Kind())
	}
	t := TypeOf(x)
	size := t.Len()
	dc := New(ArrayOf(size, t.Elem())).Elem()
	for i := 0; i < size; i++ {
		item, err := _anything(v.Index(i).Interface(), ptrs)
		if err != nil {
			return nil, fmt.Errorf("failed to clone array item at index %v: %v", i, err)
		}
		dc.Index(i).Set(ValueOf(item))
	}
	return dc.Interface(), nil
}
