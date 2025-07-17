package datetime_test

import (
	"testing"
	"time"

	"github.com/jummyliu/pkg/datetime"
)

func TestGetMonthRange(t *testing.T) {
	now := time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC)
	start, end := datetime.GetMonthRange(now)

	targetStart := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	targetEnd := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	if !(start.Equal(targetStart) && end.Equal(targetEnd)) {
		t.Fatalf("Test failed, need (%s-%s) but got (%s-%s)", targetStart, targetEnd, start, end)
	}
}

func TestGetWeekRange(t *testing.T) {
	now := time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC)
	start, end := datetime.GetWeekRange(now, time.Monday)

	targetStart := time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC)
	targetEnd := time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC)

	if !(start.Equal(targetStart) && end.Equal(targetEnd)) {
		t.Fatalf("Test failed, need (%s-%s) but got (%s-%s)", targetStart, targetEnd, start, end)
	}
}

func TestGetYearRange(t *testing.T) {
	now := time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC)
	start, end := datetime.GetYearRange(now)

	targetStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	targetEnd := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	if !(start.Equal(targetStart) && end.Equal(targetEnd)) {
		t.Fatalf("Test failed, need (%s-%s) but got (%s-%s)", targetStart, targetEnd, start, end)
	}
}

func TestGetDayRange(t *testing.T) {
	now := time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC)
	start, end := datetime.GetDayRange(now)

	targetStart := time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC)
	targetEnd := time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC)

	if !(start.Equal(targetStart) && end.Equal(targetEnd)) {
		t.Fatalf("Test failed, need (%s-%s) but got (%s-%s)", targetStart, targetEnd, start, end)
	}
}
