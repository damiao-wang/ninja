package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"ninja/base/misc/errors"
	"ninja/base/misc/stack"

	"github.com/Sirupsen/logrus"
)

var std = NewEx(1)

var errHandler func(err error, depth int)

type Entry struct {
	*logrus.Entry
	depth int
}

func New() *Entry {
	return NewEx(0)
}

func SetErrorHandler(f func(err error, depth int)) {
	errHandler = f
}

func NewEx(depth int) *Entry {
	logger := logrus.New()
	disableColor := os.Getenv("TERM") == "dumb"
	if !disableColor {
		disableColor = os.Getenv("NOCOLOR") != ""
	}
	logger.Formatter = &logrus.TextFormatter{
		DisableColors: disableColor,
	}
	return &Entry{
		Entry: logrus.NewEntry(logger),
		depth: depth,
	}
}

func (e *Entry) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		e.Entry.Logger.Level = logrus.DebugLevel
	case "info":
		e.Entry.Logger.Level = logrus.InfoLevel
	case "warn":
		e.Entry.Logger.Level = logrus.WarnLevel
	case "error":
		e.Entry.Logger.Level = logrus.ErrorLevel
	case "fatal":
		e.Entry.Logger.Level = logrus.FatalLevel
	}
}

func (e *Entry) Printf(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Printf(format, args...)
}

func (e *Entry) Println(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Println(args...)
}

func (e *Entry) Debug(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Debug(args...)
}

func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Debugf(format, args...)
}

func (e *Entry) Fatal(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Fatal(args...)
}

func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Fatalf(format, args...)
}

func (e *Entry) Info(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Info(args...)
}

func (e *Entry) Infof(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Infof(format, args...)
}

func (e *Entry) Warn(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Warn(args...)
}

func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Warnf(format, args...)
}

func (e *Entry) WithField(key string, value interface{}) *Entry {
	return &Entry{e.Entry.WithField(key, value), e.depth}
}

func (e *Entry) Fields(fields map[string]interface{}) *Entry {
	return &Entry{e.Entry.WithFields(fields), e.depth}
}

func (e *Entry) Field(key string, value interface{}) *Entry {
	return &Entry{e.Entry.WithField(key, value), e.depth}
}

func (e *Entry) Error(err error, args ...interface{}) {
	if e.depth >= 0 {
		err = errors.WrapEx(e.depth+2, err, nil)
	}
	entry := e.Entry
	if es, ok := err.(Fielder); ok {
		if fields := es.Fields(); fields != nil {
			entry = e.Entry.WithFields(map[string]interface{}(fields))
		}
	}
	if errHandler != nil {
		errHandler(err, e.depth)
	}
	if len(args) == 0 {
		entry.Error(err.Error())
	} else {
		entry.WithError(err).Error(args...)
	}
}

func (e *Entry) Errorf(format string, args ...interface{}) {
	err := fmt.Errorf(format, args...)
	if errHandler != nil {
		errHandler(err, e.depth)
	}
	e.Entry.WithField("call", stack.String(e.depth+1)).Error(err)
}

func (e *Entry) ErrorMsg(args ...interface{}) {
	e.Entry.WithField("call", stack.String(e.depth+1)).Error(args...)
}

// -----------------------------------------------------------------------------

type Fielder interface {
	Fields() stack.Fields
}

func Field(key string, val interface{}) *Entry {
	return New().Field(key, val)
}

func Fields(fields map[string]interface{}) *Entry {
	return New().Fields(fields)
}

func SetLevel(level string) {
	std.SetLevel(level)
}

func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	std.Fatalln(args...)
}

func Error(err error, args ...interface{}) {
	std.Error(err, args...)
}

func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

func ErrorMsg(args ...interface{}) {
	std.ErrorMsg(args...)
}

func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}

func Println(args ...interface{}) {
	std.Println(args...)
}

func Info(args ...interface{}) {
	std.Info(args...)
}

func Infof(layout string, args ...interface{}) {
	std.Infof(layout, args...)
}

func Warn(args ...interface{}) {
	std.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

func Struct(obj ...interface{}) {
	for _, o := range obj {
		fmt.Printf("%T(%+v)\n", o, o)
	}
}

func JSON(obj interface{}) {
	ret, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(ret))
}
