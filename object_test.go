package tox

import (
	"github.com/davecgh/go-spew/spew"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestReportSuite(t *testing.T) {
	suite.Run(t, new(ReportSuite))
}

type ReportSuite struct {
	suite.Suite
}

func (s *ReportSuite) SetupSuite() {

}

func (s *ReportSuite) TestSimpleFlatten() {

	o := Object{"a": "abc", "b": 456, "c": Object{"d": 123, "e": Object{"f": 456}}}
	s.Equal(Object{"a": "abc", "b": 456, "c.d": 123, "c.e.f": 456}, o.Flatten("."))
}

func (s *ReportSuite) TestArrayFlatten() {

	o := Object{"a": "abc", "b": 456, "c": Object{"d": 123, "e": Object{"f": 456}, "g": []string{"h", "i", "j"}}}
	s.Equal(Object{"a": "abc", "b": 456, "c.d": 123, "c.e.f": 456, "c.g[0]": "h", "c.g[1]": "i", "c.g[2]": "j"}, o.Flatten("."))
}

func (s *ReportSuite) TestDiffMod() {

	old := Object{"a": "abc", "b": 456, "c": Object{"d": 123, "e": Object{"f": 456}}}
	new2 := Object{"a": "abc", "b": 456, "c": Object{"d": 555, "e": Object{"f": 456}}}

	s.Equal(ObjectDiff{Added: nil, Deleted: nil, Modified: map[string]FieldDiff{"c/d": {Old: 123, New: 555}}}, old.Diff(new2))
}

func (s *ReportSuite) TestDiff() {

	old := Object{"a": "abc", "b": 456, "c": Object{"d": 123, "e": Object{"f": 456}}}
	new2 := Object{"a": "abc", "c": Object{"d": 555, "e": Object{"f": 456, "g": 789}}}

	s.Equal(ObjectDiff{Added: Object{"c/e/g": 789}, Deleted: Object{"b": 456}, Modified: map[string]FieldDiff{"c/d": {Old: 123, New: 555}}}, old.Diff(new2))
}

type Foo struct {
	A string  `json:"a,omitempty"`
	B float64 `json:"b,omitempty"`
}

type FooUnexported struct {
	A string `json:"a,omitempty"`
	b float64
	C Foo `json:"c,omitempty"`
}

func (s *ReportSuite) TestNaN() {
	old := Object{"a": "abc", "b": math.NaN(), "c": Object{"d": 123, "e": math.NaN()}}
	old.RemoveNaN()
	s.Equal("{\"a\":\"abc\",\"c\":{\"d\":123}}", old.JsonString(false))

	old = Object{"a": "abc", "b": math.NaN(), "c": &Foo{A: "abc", B: math.NaN()}}
	old.RemoveNaN()
	s.Equal("{\"a\":\"abc\",\"c\":{\"a\":\"abc\"}}", old.JsonString(false))

	old = Object{"a": "abc", "b": math.NaN(), "c": Object{"d": &Foo{A: "abc", B: math.NaN()}}}
	old.RemoveNaN()
	s.Equal("{\"a\":\"abc\",\"c\":{\"d\":{\"a\":\"abc\"}}}", old.JsonString(false))
}

func (s *ReportSuite) TestBinary() {
	old := Object{"a": "abc", "b": []byte{0x01, 0x02, 0x03, 0x04}}
	n := NewObject(map[string]any{"a": "abc", "b": []byte{0x01, 0x02, 0x03, 0x04}})
	n.Equals(old)
	spew.Dump(old, n)
	s.Equal(true, n.Equals(old))
}

func (s *ReportSuite) TestStructs() {
	old := NewObject(map[string]interface{}{"a": "abc", "b": 123})
	old.Set("c", Foo{A: "abc", B: 456})
	s.Equal("{\"a\":\"abc\",\"b\":123,\"c\":{\"a\":\"abc\",\"b\":456}}", old.JsonString(false))

	old = NewObject(map[string]interface{}{"a": "abc", "b": 123})
	old.Set("c", FooUnexported{A: "abc", b: 456})
	s.Equal("{\"a\":\"abc\",\"b\":123,\"c\":{\"a\":\"abc\",\"c\":{}}}", old.JsonString(false))

	old = Object{"a": "abc", "b": 123, "c": FooUnexported{A: "abc", b: 456}}
	s.Equal("{\"a\":\"abc\",\"b\":123,\"c\":{\"a\":\"abc\",\"c\":{}}}", old.JsonString(false))

	old = Object{"a": "abc", "b": 123, "c": &FooUnexported{A: "abc", b: 456}}
	s.Equal("{\"a\":\"abc\",\"b\":123,\"c\":{\"a\":\"abc\",\"c\":{}}}", old.JsonString(false))

	old = Object{"a": "abc", "b": 123, "c": &FooUnexported{A: "abc", b: 456, C: Foo{}}}
	s.Equal("{\"a\":\"abc\",\"b\":123,\"c\":{\"a\":\"abc\",\"c\":{}}}", old.JsonString(false))

	old = Object{"a": "abc", "b": 123, "c": &FooUnexported{A: "abc", b: 456, C: Foo{A: "abc", B: 456}}}
	s.Equal("{\"a\":\"abc\",\"b\":123,\"c\":{\"a\":\"abc\",\"c\":{\"a\":\"abc\",\"b\":456}}}", old.JsonString(false))

	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	old = Object{"a": "abc", "b": now}
	old.RemoveNaN()
	s.Equal("{\"a\":\"abc\",\"b\":\"2023-01-01T00:00:00Z\"}", old.JsonString(false))
}

func (s *ReportSuite) TestGetFloat64Array() {
	o := Object{
		"a": []float64{1.1, 2.2},
		"b": []int{1, 2},
		"c": "not an array",
	}

	s.Equal([]float64{1.1, 2.2}, o.GetFloat64Array("a", nil))
	s.Equal([]float64{1.0, 2.0}, o.GetFloat64Array("b", nil))
	res := o.GetFloat64Array("c", nil)
	s.Len(res, 1)
	s.True(math.IsNaN(res[0]), "Should return NaN for invalid float string in array conversion")
	s.Equal([]float64{9.9}, o.GetFloat64Array("d", []float64{9.9}))
	s.Nil(o.GetFloat64Array("d", nil))
}
