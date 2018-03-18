package errors

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"unsafe"

	"ninja/base/misc/stack"
)

type Error interface {
	error
	Stack() []string
}

var _ coder = &stackError{}

type stackError struct {
	stack    []*stack.Frame
	fields   stack.Fields
	tags     map[string]string
	errCode  *int32
	isIgnore bool
	prefix   []string
	error
}

func (s *stackError) Copy() *stackError {
	stack := make([]*stack.Frame, len(s.stack))
	copy(stack, s.stack)
	tags := make(map[string]string, len(s.tags))
	for k, v := range s.tags {
		tags[k] = v
	}
	prefix := make([]string, len(s.prefix))
	copy(prefix, s.prefix)

	return &stackError{
		stack:    stack,
		fields:   s.fields.Clone(),
		tags:     tags,
		errCode:  s.errCode,
		isIgnore: s.isIgnore,
		prefix:   prefix,
		error:    s.error,
	}
}

func (s *stackError) Code() int32 {
	if s.errCode != nil {
		return *s.errCode
	}
	if errCoder, ok := s.error.(coder); ok {
		return errCoder.Code()
	}
	return int32(ErrUnknown)
}

func (s *stackError) Error() string {
	if len(s.prefix) > 0 {
		return strings.Join(s.prefix, ".") + ": " + s.error.Error()
	}
	return s.error.Error()
}

func (s *stackError) getFrameString() string {
	buf := bytes.NewBuffer(nil)
	for _, item := range s.stack {
		buf.WriteString(item.String())
		buf.WriteString(";")
	}
	return buf.String()
}

func (s *stackError) Fields() stack.Fields {
	stack := s.getFrameString()
	if stack == "" {
		return s.fields
	}
	if s.fields == nil {
		return map[string]interface{}{
			"stack": s.getFrameString(),
		}
	}
	s.fields["stack"] = s.getFrameString()
	return s.fields
}

func (s *stackError) Stack() []*stack.Frame {
	return s.stack
}

func New(msg string) error {
	return WrapEx(2, fmt.Errorf(msg), nil)
}

func Fmt(layout string, args ...interface{}) error {
	return WrapEx(2, fmt.Errorf(layout, args...), nil)
}

func NewBy(code ErrorCode, layout string, args ...interface{}) error {
	err := WrapEx(2, fmt.Errorf(layout, args...), nil)
	err = SetCode(err, code)
	return err
}

type FormatError struct {
	ptr    uintptr
	msg    string
	fields stack.Fields
}

func (f *FormatError) Error() string {
	return f.msg
}

func (f *FormatError) Fields(fields map[string]interface{}) *FormatError {
	f.fields = fields
	return f
}

func (f FormatError) FormatEx(depth int, args ...interface{}) error {
	f.msg = fmt.Sprintf(f.msg, args...)
	return WrapEx(2+depth, &f, f.fields)
}

func (f FormatError) Format(args ...interface{}) error {
	return f.FormatEx(0, args...)
}

func (f FormatError) Panic(args ...interface{}) error {
	panic(f.Format(args...).Error())
}

func Define(layout string) *FormatError {
	fe := &FormatError{msg: layout}
	fe.ptr = uintptr(unsafe.Pointer(fe))
	return fe
}

func Trace(err error) error {
	return WrapEx(2, err, nil)
}

func TraceWithCode(err error, code ErrorCode) error {
	err = SetCode(err, code)
	return WrapEx(2, err, nil)
}

func TracePrefix(err error, prefix string) error {
	ret := WrapEx(2, err, nil)
	se := ret.(*stackError)
	se.prefix = append(se.prefix, prefix)
	return ret
}

func TraceWithFields(err error, fields stack.Fields) error {
	return WrapEx(2, err, fields)
}

func TraceWithField(err error, fields ...interface{}) error {
	return WrapEx(2, err, stack.Field(fields...))
}

func WrapEx(depth int, err error, fields stack.Fields) error {
	if err == nil {
		return nil
	}
	frame := stack.GetFrame(depth)

	s, ok := err.(*stackError)
	if !ok {
		return &stackError{[]*stack.Frame{frame}, fields.Clone(),
			nil, nil, false, nil, err}
	}

	s = s.Copy()
	s.stack = append(s.stack, frame)
	if s.fields == nil {
		s.fields = fields.Clone()
	} else {
		for k, v := range fields {
			s.fields[k] = v
		}
	}
	return s
}

func Replace(err, newerr error) error {
	if err == nil {
		return nil
	}
	if s, ok := err.(*stackError); ok {
		s.error = newerr
		return s
	}
	return newerr
}

func AddField(err error, fields stack.Fields) error {
	if err == nil {
		return nil
	}
	s, ok := err.(*stackError)
	if !ok {
		return &stackError{[]*stack.Frame{}, fields.Clone(), nil, nil, false, nil, err}
	}
	if s.fields == nil {
		s.fields = fields.Clone()
	} else {
		for k, v := range fields {
			s.fields[k] = v
		}
	}
	return s
}

func AddStack(err error, stack []*stack.Frame) error {
	if err == nil {
		return nil
	}
	s, ok := err.(*stackError)
	if !ok {
		return &stackError{stack: stack, error: err}
	}
	s.stack = append(s.stack, stack...)
	return s
}

func Cause(err error) error {
	if e, ok := err.(*stackError); ok {
		return e.error
	}
	return err
}

func Is(err1, err2 error) bool {
	if e, ok := err1.(*stackError); ok {
		err1 = e.error
	}
	if e, ok := err2.(*stackError); ok {
		err2 = e.error
	}
	if e1, ok := err1.(*FormatError); ok {
		if e2, ok := err2.(*FormatError); ok {
			return e1.ptr == e2.ptr
		}
	}
	return err1 == err2
}

type fielder struct {
	data map[string]interface{}
}

func Stack(err error) []*stack.Frame {
	if se, ok := err.(*stackError); ok {
		return se.Stack()
	}
	return nil
}

func Fields(err error) stack.Fields {
	if se, ok := err.(*stackError); ok {
		return se.Fields()
	}
	return nil
}

// ---------------------------------------------------------------------------

type coder interface {
	Code() int32
}

func IgnoreIf(err error, code ErrorCode) error {
	if Code(err) == code {
		return SetIgnore(err)
	}
	return err
}

func CodeMatch(err error, code ErrorCode) bool {
	return Code(err) == code
}

func Code(err error) ErrorCode {
	if err == nil {
		return OK
	}
	// if db.IsMgoNotFound(Cause(err)) {
	// 	return ErrNotFound
	// }
	if errCoder, ok := Cause(err).(coder); ok {
		return ErrorCode(errCoder.Code())
	}
	if errCoder, ok := err.(coder); ok {
		return ErrorCode(errCoder.Code())
	}

	return ErrUnknown
}

func SetCode(err error, code ErrorCode) error {
	if err == nil {
		return nil
	}
	se, ok := err.(*stackError)
	if !ok {
		se = &stackError{nil, make(stack.Fields), nil, nil, false, nil, err}
	}

	se.errCode = (*int32)(&code)
	return se
}

type Tags map[string]string

func (t Tags) Clone() Tags {
	n := make(Tags, len(t))
	for k, v := range t {
		n[k] = v
	}
	return n
}

func AddTags(err error, tags Tags) error {
	se, ok := err.(*stackError)
	if !ok {
		return &stackError{
			nil, nil, tags.Clone(), nil, false, nil, err,
		}
	}
	if se.tags == nil {
		se.tags = tags.Clone()
		return se
	}
	for k, v := range tags {
		se.tags[k] = v
	}
	return se
}

func SetIgnore(err error) error {
	if err == nil {
		return nil
	}
	se, ok := err.(*stackError)
	if !ok {
		se = &stackError{
			nil, nil, nil, nil, false, nil, err,
		}
	}
	se.isIgnore = true
	return se
}

func IsIgnore(err error) bool {
	if err == nil {
		return false
	}
	se, ok := err.(*stackError)
	if !ok {
		return false
	}
	return se.isIgnore
}

func GetTags(err error) Tags {
	if err == nil {
		return nil
	}
	se, ok := err.(*stackError)
	if ok {
		return se.tags
	}
	return nil
}

func FrameName(err error) string {
	frames := Stack(err)
	name := ""
	if len(frames) > 0 {
		name = path.Base(frames[0].Frame.Module) + "." + frames[0].Frame.Function
	}
	return name
}
