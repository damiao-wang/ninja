package route

import (
	"log"
	"testing"
)

func BenchmarkGetRoute(b *testing.B) {
	b.StopTimer()
	item := &TblRouteCfg{
		ProdCd:    "123456",
		CardClass: "1",
	}
	pool := GetPool("127.0.0.1:6379")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if _, err := pool.GetRoute(item); err != nil {
			log.Fatal(err)
		}
	}
}
