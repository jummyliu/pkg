package types

import (
	"testing"
)

func TestSliceScan(t *testing.T) {
	from := `["1","2","3"]`
	length := 3
	result1 := "2"
	s := Slice[string]{}
	err := s.Scan([]byte(from))
	if err != nil {
		t.Fatalf("Slice.Scan failure: %s", err)
	}
	if length != len(s) {
		t.Fatalf("Slice.Scan failure: %s", "length is not equal")
	}
	if result1 != s[1] {
		t.Fatalf("Slice.Scan(s[1]) need %s but got %s", result1, s[1])
	}
}

func TestSliceValue(t *testing.T) {
	result := `[1,"2",3]`
	m := Slice[any]{}
	m = append(m, 1, "2", 3)
	v, err := m.Value()
	if err != nil {
		t.Fatalf("Slice.Value failure: %s", err)
	}
	val, ok := v.([]byte)
	if !ok {
		t.Fatalf("Slice.Value failure: %s", "cover to []byte failure")
	}
	if result != string(val) {
		t.Fatalf("Slice.Value need %s but got %s", result, string(val))
	}
}
