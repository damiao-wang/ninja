package context

import (
	"testing"

	"ninja/base/misc/errors"
	"ninja/base/misc/stack"
)

var (
	ErrNotFound = errors.Define("%v is not defined")
)

func func2(ct T, name string) error {
	ct.SetFields(stack.Field("name", name))
	return ct.Format(ErrNotFound, name)
}

func func1(ctx T) error {
	if err := func2(ctx, "1"); err != nil {
		return ctx.Trace(err)
	}
	return nil
}

func TestError(t *testing.T) {
	ctx := Dump()
	ctx.SetFields(stack.Field("test", 1))
	if err := func1(ctx); err != nil {
		ctx.LogError(err)
	}
	ctx.LogInfo("hello?")

	LogInfo("hello", "a", 1)
}
