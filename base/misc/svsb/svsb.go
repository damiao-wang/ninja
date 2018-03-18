package svsb

import (
	"encoding/json"
	"reflect"
	"unsafe"
)

// 直接加引用而不是复制
// 仅适用于确认不会通过 []byte 修改内容, 导致违背string不可变的原则
// https://github.com/golang/go/commit/71be0138421012d04e06991d37d19c9f5b1fa02b
// https://github.com/golang/go/commit/e6fac08146df323eb95f46508bef937cdfb802fd

func String(b []byte) string {
	sliceh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	stringh := &reflect.StringHeader{
		Data: sliceh.Data,
		Len:  sliceh.Len,
	}
	s := (*string)(unsafe.Pointer(stringh))
	return *s
}

func Bytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh := &reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	b := (*[]byte)(unsafe.Pointer(sh))
	return *b
}

func JsonMarshal(o interface{}) (string, error) {
	ret, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return String(ret), nil
}

func JsonUnmarshal(s string, obj interface{}) error {
	return json.Unmarshal(Bytes(s), obj)
}
