package utils

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestUUID(t *testing.T) {
	str := UUID()
	expectVal := 36
	actualVal := len(str)
	if expectVal != actualVal {
		t.Fatalf("UUID's len need %d but got %d", expectVal, actualVal)
	}
}

func TestExecutablePath(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	expectVal := filepath.Dir(file)
	actualVal := GetExecutablePath()
	if expectVal != actualVal {
		t.Fatalf("executable path need %s but got %s", expectVal, actualVal)
	}
}

func FuzzRandomStrWithSource(f *testing.F) {
	testcases := []string{
		"hello world",
		" ",
		"!12345",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, str string) {
		l := 16
		target := RandomStrWithSource(l, str)
		if len(target) != l && len(str) != 0 {
			t.Fatalf("random str's len need %d bug got %d", l, len(target))
		}
	})
}
