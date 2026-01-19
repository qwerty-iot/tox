package tox

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name    string    `json:"name"`
	Age     int       `json:"age,omitempty"`
	Created time.Time `json:"created"`
	Hidden  string    `json:"-"`
	private string
}

func TestStructToObject(t *testing.T) {
	now := time.Now()
	ts := TestStruct{
		Name:    "John",
		Age:     30,
		Created: now,
		Hidden:  "secret",
		private: "shhh",
	}

	obj := structToObject(ts)
	assert.NotNil(t, obj)
	assert.Equal(t, "John", obj["name"])
	assert.Equal(t, 30, obj["age"])

	// In the improved implementation, we want time.Time to remain time.Time
	created, ok := obj["created"].(time.Time)
	assert.True(t, ok, "Expected time to be time.Time")
	assert.True(t, now.Equal(created))

	assert.Nil(t, obj["Hidden"])
	assert.Nil(t, obj["private"])
}

func TestStructToObjectOmitempty(t *testing.T) {
	ts := TestStruct{
		Name: "John",
		Age:  0, // should be omitted
	}

	obj := structToObject(ts)
	assert.NotNil(t, obj)
	assert.Equal(t, "John", obj["name"])
	_, exists := obj["age"]
	assert.False(t, exists, "age should be omitted")
}

type NestedStruct struct {
	Title string     `json:"title"`
	Data  TestStruct `json:"data"`
}

func TestNestedStructToObject(t *testing.T) {
	now := time.Now()
	ns := NestedStruct{
		Title: "Report",
		Data: TestStruct{
			Name:    "Jane",
			Created: now,
		},
	}

	obj := structToObject(ns)
	assert.NotNil(t, obj)
	assert.Equal(t, "Report", obj["title"])

	data, ok := obj["data"].(Object)
	assert.True(t, ok, "Expected data to be Object")
	assert.Equal(t, "Jane", data["name"])
	assert.True(t, now.Equal(data["created"].(time.Time)))
}

func TestSliceOfStructs(t *testing.T) {
	slice := []TestStruct{
		{Name: "A"},
		{Name: "B"},
	}

	res := structToAnything(slice)
	s, ok := res.([]any)
	assert.True(t, ok)
	assert.Equal(t, 2, len(s))
	assert.Equal(t, "A", s[0].(Object)["name"])
	assert.Equal(t, "B", s[1].(Object)["name"])
}

func TestMapToObject(t *testing.T) {
	m := map[string]TestStruct{
		"first": {Name: "A"},
	}

	res := structToAnything(m)
	om, ok := res.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "A", om["first"].(Object)["name"])
}

func TestPointerToStruct(t *testing.T) {
	ts := &TestStruct{Name: "Ptr"}
	obj := structToObject(ts)
	assert.Equal(t, "Ptr", obj["name"])

	var nilTs *TestStruct
	assert.Nil(t, structToObject(nilTs))
}
