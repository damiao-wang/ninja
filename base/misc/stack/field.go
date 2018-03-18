package stack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"ninja/base/misc/svsb"

	raven "github.com/getsentry/raven-go"
)

func PkgName(depth int) string {
	pc, _, _, ok := runtime.Caller(1 + depth)
	if !ok {
		return ""
	}
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

type Frame struct {
	Name  string
	Frame *raven.StacktraceFrame
}

func GetAllFrame(skip int) (ret []*Frame) {
	for i := skip + 1; ; i++ {
		item := GetFrame(i)
		if item == nil {
			break
		}
		if item.Frame == nil {
			break
		}
		if strings.HasPrefix(item.Name, "runtime.") {
			continue
		}
		ret = append(ret, item)
	}
	return ret
}

func GetFrame(skip int) *Frame {
	pc, file, line, ok := runtime.Caller(1 + skip)
	if !ok {
		return nil
	}
	return &Frame{
		Name: fmt.Sprintf("%v:%v",
			filepath.Base(runtime.FuncForPC(pc).Name()), line),
		Frame: raven.NewStacktraceFrame(pc, file, line, 4, nil),
	}
}

func FramesUnmarshal(data []string) (ret []*Frame, err error) {
	if len(data) == 0 {
		return nil, nil
	}
	err = json.Unmarshal(svsb.Bytes(data[0]), &ret)
	return
}

func FramesMarshal(fs []*Frame) string {
	bty, _ := json.Marshal(fs)
	return svsb.String(bty)
}

func (f *Frame) String() string {
	return f.Name
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
