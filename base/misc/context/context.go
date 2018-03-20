package context

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"

	"ninja/base/misc/errors"
	"ninja/base/misc/log"
	"ninja/base/misc/stack"
)

func From(c context.Context) T {
	if t, ok := c.(T); ok {
		return t
	}
	if c == nil {
		c = context.TODO()
	}
	return newctx(c)
}

func Dump() T {
	return newctx(context.TODO())
}

type Context context.Context

func WithDeadline(parent T, dead time.Time) (T, CancelFunc) {
	ctx, cancel := context.WithDeadline(parent.Context, dead)
	return newctx(ctx), CancelFunc(cancel)
}

func WithCancel(parent T) (T, CancelFunc) {
	ctx, cancel := context.WithCancel(parent.Context)
	return newctx(ctx), CancelFunc(cancel)
}

type CancelFunc context.CancelFunc

type T struct {
	context.Context
	m *sync.Mutex
}

func newctx(c context.Context) T { return T{Context: c, m: new(sync.Mutex)} }

func (c T) Format(e *errors.FormatError, args ...interface{}) error {
	return e.Fields(c.Fields()).FormatEx(1, args...)
}

func (c T) Code(err error) errors.ErrorCode {
	return errors.Code(err)
}

func (c *T) Update(c2 context.Context) {
	c.Context = c2
}

func (c *T) HasValue(key string) bool {
	return c.GetValue(key) == nil
}

func (c *T) SetField(args ...interface{}) {
	c.SetFields(stack.Field(args...))
}

func (c *T) SetValue(key interface{}, value interface{}) {
	ctx := context.WithValue(c.Context, key, value)
	c.Context = ctx
}

// 每次Request共享的
type requestHeapKey struct{}
type RequestHeap struct {
	Cookie []*http.Cookie
	Start  time.Time
}

// 这个结构体确保在请求内是共享的
func (c *T) InitRequestHeap(h *RequestHeap) *RequestHeap {
	if h == nil {
		h = new(RequestHeap)
	}
	c.SetValue(requestHeapKey{}, h)
	return h
}

func (c *T) GetRequestHeap() *RequestHeap {
	return c.GetValue(requestHeapKey{}).(*RequestHeap)
}

func (c *T) SetResponseCookie(cks ...*http.Cookie) {
	ck := c.GetResponseCookie()
	heap := c.GetRequestHeap()
	if len(ck) == 0 {
		heap.Cookie = cks
	} else {
		heap.Cookie = append(heap.Cookie, cks...)
	}
}

type ctxHttpReq struct{}

func (c *T) SetRequest(req *http.Request) {
	c.SetValue(ctxHttpReq{}, req)
}

func (c *T) GetCookie(key string) (*http.Cookie, error) {
	req := c.GetRequest()
	if req != nil {
		return req.Cookie(key)
	}

	md := c.GetInMetadata()
	grpcCookies := md["cookie"]

	for _, it := range grpcCookies {
		sp := strings.Split(it, ";")
		for _, item := range sp {
			item = strings.TrimSpace(item)
			kv := strings.Split(item, "=")
			if len(kv) == 2 && strings.EqualFold(kv[0], key) {
				return &http.Cookie{
					Value: kv[1],
					Name:  kv[0],
				}, nil
			}
		}
	}

	return nil, errors.Fmt("http.Request is not set or GRPC metadata is not set")

}

func (c *T) SingleFlight(f func()) {
	c.m.Lock()
	f()
	c.m.Unlock()
}

// 获取grpc的metadata
func (c *T) GetGrpcMetaValues(keys ...string) (map[string]string, error) {
	if len(keys) <= 0 {
		return nil, errors.Fmt("key is nil.")
	}

	md := c.GetInMetadata()
	ret := make(map[string]string)
	for _, key := range keys {
		grpcMetaValue := md[key]
		if len(grpcMetaValue) > 0 {
			ret[key] = grpcMetaValue[0]
		}
	}
	if len(ret) > 0 {
		return ret, nil
	}

	return nil, errors.Fmt("GRPC metadata is not set")
}

// 获取grpc的metadata
func (c *T) GetGrpcMetaValue(key string) string {
	var val string
	md := c.GetInMetadata()
	grpcMetaValue := md[key]
	if len(grpcMetaValue) > 0 {
		val = grpcMetaValue[0]
	}
	return val
}

func (c *T) SetMeta(key string, value string) {
	c.GetOutMetadata()[key] = []string{value}
	c.GetInMetadata()[key] = []string{value}
}

func (c *T) SetDebug(info string) {
	c.SetMeta("ezbuy-debug", info)
}

func (c *T) IsDebugFor(name string) bool {
	sp := strings.Split(c.GetDebug(), ",")
	for _, item := range sp {
		if strings.EqualFold(name, strings.TrimSpace(item)) {
			return true
		}
	}
	return false
}

func (c *T) GetDebug() string {
	req := c.GetRequest()
	if req != nil {
		info := req.Header.Get("ezbuy-debug")
		if info != "" {
			return info
		}
	}
	return c.GetGrpcMetaValue("ezbuy-debug")
}

func (c *T) GetRequest() *http.Request {
	req := c.GetValue(ctxHttpReq{})
	if req == nil {
		return nil
	}
	if r, ok := req.(*http.Request); ok {
		return r
	}
	return nil
}

func (c *T) GetValue(key interface{}) interface{} {
	return c.Context.Value(key)
}

func (c *T) GetResponseCookie() []*http.Cookie {
	return c.GetRequestHeap().Cookie
}

func (c *T) SetFields(f stack.Fields) {
	f.Merge(c.Fields())
	ctx := context.WithValue(c.Context, "fields", f)
	c.Context = ctx
}

func (c T) LogError(err error) {
	c.LogErrorEx(1, err)
}

func (c T) LogErrorf(layout string, args ...interface{}) {
	c.LogErrorEx(1, errors.Fmt(layout, args...))
}

func (c T) LogErrorEx(depth int, err error) {
	if depth >= 0 {
		depth += 1
	}
	log.NewEx(depth).Error(errors.AddField(err, c.Fields()))
}

func (c *T) GetHeader(key string) string {
	md := c.GetInMetadata()
	value := md[key]
	if len(value) > 0 {
		return value[0]
	}
	return ""
}

func (c *T) GetOutMetadata() metadata.MD {
	md, ok := metadata.FromOutgoingContext(c.Context)
	if !ok {
		md = metadata.MD{}
		c.Context = metadata.NewOutgoingContext(c.Context, md)
	}
	return md
}

func (c *T) GetInMetadata() metadata.MD {
	md, ok := metadata.FromIncomingContext(c.Context)
	if !ok {
		md = metadata.MD{}
		c.Context = metadata.NewIncomingContext(c.Context, md)
	}
	return md
}

func (c *T) SendHeader(key string, value string) {
	md := c.GetOutMetadata()
	md[key] = []string{value}
}

func (c T) LogInfoEx(depth int, msg string, args ...interface{}) {
	if depth >= 0 {
		depth += 1
	}
	log.NewEx(depth).Fields(c.Fields()).Info(fmt.Sprintf(msg, args...))
}

func (c T) LogInfof(msg string, args ...interface{}) {
	c.LogInfoEx(1, msg, args...)
}

func (c T) LogInfo(msg string, args ...interface{}) {
	c.LogInfoEx(1, msg, args...)
}

func (c T) Fields() stack.Fields {
	fs, _ := c.Context.Value("fields").(stack.Fields)
	return fs
}

func (c T) Cause(err error) error {
	return errors.Cause(err)
}

func (c T) CodeErrorf(code errors.ErrorCode, layout string, args ...interface{}) error {
	err := errors.SetCode(errors.Fmt(layout, args...), code)
	return errors.WrapEx(2, err, c.Fields())
}

func (c T) TraceCode(err error, code errors.ErrorCode) error {
	err = errors.SetCode(err, code)
	return errors.WrapEx(2, err, c.Fields())
}

func (c T) Trace(err error) error {
	return errors.WrapEx(2, err, c.Fields())
}

func (c T) Errorf(msg string, args ...interface{}) error {
	return errors.WrapEx(2, fmt.Errorf(msg, args...), c.Fields())
}

func (c T) TraceWithField(err error, fields ...interface{}) error {
	return errors.WrapEx(2, err, stack.Field(fields...))
}

func LogError(err error, fields ...interface{}) {
	logger.Error(errors.AddField(err, stack.Field(fields...)))
}

func LogJSON(obj interface{}) {
	data, _ := json.MarshalIndent(obj, "", "  ")
	println(string(data))
}

func LogObj(obj interface{}) {
	logger.Infof("%+v", obj)
}

func LogInfof(msg string, args ...interface{}) {
	logger.Info(fmt.Sprintf(msg, args...))
}

func LogInfo(msg string, fields ...interface{}) {
	if len(fields)%2 != 0 {
		logger.Info(fmt.Sprintf(msg, fields...))
		return
	}
	logger.Fields(stack.Field(fields...)).Info(msg)
}

var logger = log.NewEx(1)

func TraceWithField(err error, fields ...interface{}) error {
	return errors.WrapEx(2, err, stack.Field(fields...))
}
