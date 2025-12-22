package tox

import (
	"testing"
	"time"
)

func TestDeepcopyTime(t *testing.T) {
	now := time.Now()
	ptrNow := &now
	m := map[string]any{
		"time":     now,
		"ptr_time": ptrNow,
		"nil_time": (*time.Time)(nil),
	}

	cp, err := Deepcopy(m)
	if err != nil {
		t.Fatalf("Deepcopy failed: %v", err)
	}

	m2 := any(cp).(map[string]any)
	now2 := m2["time"].(time.Time)

	if !now.Equal(now2) {
		t.Errorf("expected %v, got %v", now, now2)
	}

	ptrNow2 := m2["ptr_time"].(*time.Time)
	if ptrNow2 == ptrNow {
		t.Errorf("expected different pointer for ptr_time")
	}
	if !ptrNow.Equal(*ptrNow2) {
		t.Errorf("expected %v, got %v", *ptrNow, *ptrNow2)
	}

	if m2["nil_time"].(*time.Time) != nil {
		t.Errorf("expected nil for nil_time")
	}
}
