package htl

import "testing"

func TestLogBeforeStart(t *testing.T) {
	if Log() != nil {
		t.Fatal("Log antes de Start deveria ser nil")
	}
}
