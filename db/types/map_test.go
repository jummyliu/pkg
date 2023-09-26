package types

import (
	"testing"
)

func TestMapScan(t *testing.T) {
	from := `{"hello":"world"}`
	result := "world"
	m := Map[string, string]{}
	err := m.Scan([]byte(from))
	if err != nil {
		t.Fatalf("Map.Scan failure: %s", err)
	}
	if result != m["hello"] {
		t.Fatalf("Map.Scan(m['hello']) need %s but got %s", result, m["hello"])
	}
}

func TestMapValue(t *testing.T) {
	result := `{"hello":"world"}`
	m := Map[string, any]{}
	m["hello"] = "world"
	v, err := m.Value()
	if err != nil {
		t.Fatalf("Map.Value failure: %s", err)
	}
	val, ok := v.([]byte)
	if !ok {
		t.Fatalf("Map.Value failure: %s", "cover to []byte failure")
	}
	if result != string(val) {
		t.Fatalf("Map.Value need %s but got %s", result, string(val))
	}
}

func TestMustGetValue(t *testing.T) {
	m := map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"hello":   "world",
				"int":     10,
				"float64": 3.14,
			},
		},
	}
	strResult := "world"
	intResult := 10
	floatResult := 3.14
	if val := MustGetValue[string](m, "level1.level2.hello"); strResult != val {
		t.Fatalf("m[level1][level2][hello] need %s but got %s", strResult, val)
	}
	if val := MustGetValue[int](m, "level1.level2.int"); intResult != val {
		t.Fatalf("m[level1][level2][int] need %d but got %d", intResult, val)
	}
	if val := MustGetValue[float64](m, "level1.level2.float64"); floatResult != val {
		t.Fatalf("m[level1][level2][hello] need %f but got %f", floatResult, val)
	}
}
