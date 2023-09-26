package types

import (
	"testing"
	"time"

	"github.com/jummyliu/pkg/datetime"
)

func TestTimeScan(t *testing.T) {
	from := `"2023-09-26 14:42:41"`
	result := datetime.ParseDate(from)
	ti := Time{}
	err := ti.Scan([]byte(from))
	if err != nil {
		t.Fatalf("Time.Scan failure: %s", err)
	}
	if !result.Equal(time.Time(ti)) {
		t.Fatalf("Time.Scan(s[1]) need %s but got %s", result, time.Time(ti))
	}
}

func TestTimeValue(t *testing.T) {
	from := `2023-09-26 14:42:41`
	result := `2023-09-26 14:42:41 +0800 CST`
	ti := Time(datetime.ParseDate(from))
	v, err := ti.Value()
	if err != nil {
		t.Fatalf("Time.Value failure: %s", err)
	}
	val, ok := v.(time.Time)
	if !ok {
		t.Fatalf("Time.Value failure: %s", "cover to time.Time failure")
	}
	if result != val.String() {
		t.Fatalf("Time.Value need %s but got %s", result, val.String())
	}
}
