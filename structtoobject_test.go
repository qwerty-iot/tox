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

type mockObjectID [12]byte

func (id *mockObjectID) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "\"0123456789abcdef01234567\"" {
		*id = mockObjectID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	} else if s == "\"fedcba987654321012345678\"" {
		*id = mockObjectID{0xf, 0xe, 0xd, 0xc, 0xb, 0xa, 9, 8, 7, 6, 5, 4}
	}
	return nil
}

func (id mockObjectID) MarshalJSON() ([]byte, error) {
	if id == (mockObjectID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}) {
		return []byte("\"0123456789abcdef01234567\""), nil
	}
	return []byte("\"fedcba987654321012345678\""), nil
}

func (id mockObjectID) Hex() string {
	if id == (mockObjectID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}) {
		return "0123456789abcdef01234567"
	}
	return "fedcba987654321012345678"
}

type MongoStruct struct {
	Id       mockObjectID  `json:"id"`
	Name     string        `json:"name"`
	ParentId *mockObjectID `json:"parentId,omitempty"`
}

func TestMongoObjectIdHandling(t *testing.T) {
	oid := mockObjectID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	ms := MongoStruct{
		Id:   oid,
		Name: "Test",
	}

	obj := NewObject(ms)

	val := obj.Get("id")
	assert.IsType(t, "", val, "ObjectID should be encoded as a string")
	assert.Equal(t, oid.Hex(), val)

	// Test pointer to ObjectID
	parentOid := mockObjectID{0xf, 0xe, 0xd, 0xc, 0xb, 0xa, 9, 8, 7, 6, 5, 4}
	ms.ParentId = &parentOid
	obj = NewObject(ms)
	assert.Equal(t, parentOid.Hex(), obj.Get("parentId"))

	// Test ToStruct (decoding)
	var decoded MongoStruct
	obj.ToStruct(&decoded)
	assert.Equal(t, oid, decoded.Id)
	assert.Equal(t, &parentOid, decoded.ParentId)
}
