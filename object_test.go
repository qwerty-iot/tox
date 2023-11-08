package tox

import (
	"testing"

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
	s.Equal(Object{"a": "abc", "b": 456, "c.d": 123, "c.e.f": 456}, o.Flatten(""))
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
