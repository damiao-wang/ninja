package errors

import (
	"errors"
	"testing"
)

func TestTracePrefix(t *testing.T) {
	err := errors.New("hello")
	err1 := TracePrefix(err, "prefix")
	err2 := TracePrefix(err1, "prefix2")
	if err2.Error() != "prefix.prefix2: hello" {
		t.Fatal("result not expected")
	}
	if Fields(err2)["stack"] != "errors.TestTracePrefix:10;errors.TestTracePrefix:11;" {
		t.Fatal("not expected", Fields(err2))
	}
}

func TestNewBy(t *testing.T) {
	err := NewBy(ErrNotFound, "not not found")
	if !CodeMatch(err, ErrNotFound) {
		t.Fatal("result not expected")
	}
}
