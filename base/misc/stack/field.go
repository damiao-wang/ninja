package stack

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func PkgName(depth int) string {
	pc, _, _, ok := runtime.Caller(1 + depth)
	if !ok {
		return ""
	}
	fmt.Println("path: ", runtime.FuncForPC(pc).Name())
	fp := filepath.Base(runtime.FuncForPC(pc).Name())
	return strings.Split(fp, ".")[0]
}

func String(depth int) string {
	pc, _, n, ok := runtime.Caller(1 + depth)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v:%v",
		filepath.Base(runtime.FuncForPC(pc).Name()), n)
}

func All() string {
	n := bytes.NewBuffer(nil)
	for i := 1; ; i++ {
		ret := String(i)
		if ret == "" {
			break
		}
		if strings.HasPrefix(ret, "runtime.") {
			continue
		}
		n.WriteString(ret)
		n.WriteString(";")
	}
	return n.String()
}

type Fields map[string]interface{}

func Field(kvs ...interface{}) Fields {
	return make(Fields).Field(kvs...)
}

func (f Fields) Clone() Fields {
	n := make(Fields, len(f))
	for k, v := range f {
		n[k] = v
	}
	return n
}

func (f Fields) Merge(f2 Fields) {
	for k, v := range f2 {
		if _, ok := f[k]; !ok {
			f[k] = v
		}
	}
}

func (f Fields) Field(kvs ...interface{}) Fields {
	if len(kvs)&1 != 0 {
		panic("invalid kvs")
	}
	for i := 0; i < len(kvs); i += 2 {
		name := kvs[i].(string)
		f[name] = kvs[i+1]
	}
	return f
}
