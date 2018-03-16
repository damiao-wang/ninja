package stack

import (
	"path"
	"runtime"
)

func FuncName(depth int) string {
	pc, _, _, ok := runtime.Caller(1 + depth)
	if !ok {
		return ""
	}
	return runtime.FuncForPC(pc).Name()
}

func FuncNameLine(depth int) (string, int) {
	pc, file, line, ok := runtime.Caller(1 + depth)
	if !ok {
		return "", -1
	}

	return runtime.FuncForPC(pc).Name() + path.Base(file), line

}
