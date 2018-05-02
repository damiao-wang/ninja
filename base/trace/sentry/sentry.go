package sentry

import (
	"fmt"
	"strings"
	"sync"

	raven "github.com/getsentry/raven-go"

	"ninja/base/misc/context"
	"ninja/base/misc/errors"
	"ninja/base/misc/log"
	"ninja/base/misc/stack"
)

type Config struct {
	DSN string `json:"DSN"`
	sync.Once
}

const DefaultDsn string = "https://1c115aaf0b2048f485936409b03ce0f7:c6480693facf40d6b89c28711fb363b9@sentry.io/304312"

var cfg Config

func init() {
	log.SetErrorHandler(func(err error, depth int) {
		if errors.IsIgnore(err) {
			return
		}
		tags := errors.GetTags(err)
		if tags == nil {
			tags = make(map[string]string)
		}
		tags["exec"] = stack.GetServiceName()
		fields := errors.Fields(err)
		method := ""
		if m := fields["method"]; m != nil {
			method = fmt.Sprint(m)
		}
		name := errors.FrameName(err)
		if len(method) > 0 && len(name) > 0 {
			method += ":" + name
		} else {
			method += name
		}
		CaptureErrorEx(-1, method, err, tags)
	})
}

func SetDSN(dsn string) {
	raven.SetDSN(dsn)
}

func isEnable() bool {
	raven.SetDSN(DefaultDsn)
	return true
}

func last(s []string) string {
	if len(s) > 0 {
		return s[len(s)-1]
	}
	return ""
}

func CaptureErrorEx(n int, name string, err error, tags map[string]string, interfaces ...raven.Interface) string {
	if !isEnable() {
		return ""
	}
	if n >= 0 {
		err = errors.WrapEx(2+n, err, nil)
	}
	frames := errors.Stack(err)
	ravenFrames := make([]*raven.StacktraceFrame, len(frames))
	for idx, frame := range frames {
		ravenFrames[idx] = frame.Frame
		const pkg = "gitlab.1dmy.com/ezbuy"
		ravenFrames[idx].Module = strings.TrimPrefix(ravenFrames[idx].Module, pkg)
	}

	ex := raven.NewException(err, &raven.Stacktrace{ravenFrames})
	ex.Type = name
	packet := raven.NewPacket(err.Error(), ex)
	if tags == nil {
		tags = map[string]string{}
	}
	fields := errors.Fields(err)
	for k, v := range fields {
		if k == "stack" {
			continue
		}
		packet.Extra["field."+k] = v
	}
	tags["git"] = ""
	raven.DefaultClient.Capture(packet, tags)
	return ""
}

func CapturePanicEx(n int, name string, f func(), tags map[string]string, interfaces ...raven.Interface) (reterr error) {
	if !isEnable() {
		f()
		return
	}

	defer func() {
		var packet *raven.Packet
		err := recover()
		switch rval := err.(type) {
		case nil:
			return
		case error:
			rval = errors.AddStack(rval, stack.GetAllFrame(2))
			reterr = rval
			context.LogError(reterr)
			exp := raven.NewException(rval, raven.NewStacktrace(3, 3, nil))
			if name != "" {
				exp.Type = name
			}
			packet = raven.NewPacket(rval.Error(), append(interfaces, exp)...)
		default:
			rvalStr := fmt.Sprint(rval)
			packet = raven.NewPacket(rvalStr, append(interfaces, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)))...)
		}

		raven.Capture(packet, tags)
	}()

	f()
	return
}
