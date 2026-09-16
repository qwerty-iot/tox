package tox

import (
	"reflect"
	"testing"
)

func TestObjectUnmarshalTextAndBytes(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw  []byte
		kind string
	}{
		{"ASCII", []byte("device,08B9$"), "string"},
		{"ASCII whitespace", []byte("device\t08B9$\r\n"), "string"},
		{"emoji payload", []byte("device,👀,08B9$\r\n"), "string"},
		{"Unicode text", []byte("café,東京,e\u0301"), "string"},
		{"nonbreaking space", []byte("guest\u00a0network"), "string"},
		{"joined emoji", []byte("👩\u200d💻"), "string"},
		{"encoded replacement character", []byte("\ufffd"), "string"},
		{"malformed JSON text", []byte("{device:👀}"), "string"},
		{"invalid UTF-8", []byte{0xff, 0xfe}, "bytes"},
		{"truncated UTF-8", []byte{'a', 0xf0, 0x9f, 0x91}, "bytes"},
		{"overlong UTF-8", []byte{0xc0, 0xaf}, "bytes"},
		{"NUL", []byte("device\x00"), "bytes"},
		{"escape", []byte("device\x1b"), "bytes"},
		{"DEL", []byte("device\x7f"), "bytes"},
		{"Unicode control", []byte("device\u0085"), "bytes"},
		{"binary", []byte{0x01, 0x02, 0x03}, "bytes"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			obj := Object{}
			if got := obj.Unmarshal("payload", tc.raw); got != tc.kind {
				t.Fatalf("Unmarshal returned %q, want %q", got, tc.kind)
			}
			var want any = tc.raw
			if tc.kind == "string" {
				want = string(tc.raw)
			}
			if !reflect.DeepEqual(obj["payload"], want) {
				t.Fatalf("payload = %#v (%T), want %#v (%T)", obj["payload"], obj["payload"], want, want)
			}
		})
	}
}

func TestObjectUnmarshalStructuredAndEmpty(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw  []byte
		kind string
		want any
	}{
		{"object", []byte(`{"name":"👀"}`), "object", Object{"name": "👀"}},
		{"array", []byte(`[{"name":"👀"}]`), "object", []Object{{"name": "👀"}}},
		{"empty", []byte{}, "empty", nil},
		{"nil", nil, "empty", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			obj := Object{}
			if got := obj.Unmarshal("payload", tc.raw); got != tc.kind {
				t.Fatalf("Unmarshal returned %q, want %q", got, tc.kind)
			}
			if !reflect.DeepEqual(obj["payload"], tc.want) {
				t.Fatalf("payload = %#v, want %#v", obj["payload"], tc.want)
			}
		})
	}
}
