package route

import (
	"log"
	"testing"
)

func init() {
	InitRedis("127.0.0.1:6379")
}
func BenchmarkGetRoute(b *testing.B) {
	item := &TblRouteCfg{
		ProdCd:    "123456",
		CardClass: "1",
	}
	for i := 0; i < b.N; i++ {
		if _, err := GetRoute(item); err != nil {
			log.Fatal(err)
		}
	}
}
