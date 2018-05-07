package main

import (
	"fmt"
	"time"

	"ninja/base/misc/log"
	"ninja/rule2"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

var pool *redis.Pool

func main() {
	// db, _ := gorm.Open("mysql", "root:wang123456@/hyjr?charset=utf8&parseTime=True&loc=Local")
	// fmt.Println(db.AutoMigrate(&model.TblRouteCfg{}).Error)
	rule2.InitRedis("127.0.0.1:6379")
	// if err := route2.LoadRoute(db); err != nil {
	// 	log.Error(err)
	// 	return
	// }
	for {
		resp, err := rule2.GetRoute(&rule2.TblRouteCfg{
			ProdCd:    "123456",
			CardClass: "1",
		})
		if err != nil {
			log.Error(err)
		} else {
			fmt.Println("resp: ", resp)
		}
		time.Sleep(15 * time.Second)
	}
}

func testRdis() {
	newPool()
	conn1 := pool.Get()
	fmt.Printf("active:%v, ide: %v\n ", pool.Stats().ActiveCount, pool.Stats().IdleCount)
	conn1.Close()
	time.Sleep(11 * time.Second)
	fmt.Printf("active:%v, ide: %v\n ", pool.Stats().ActiveCount, pool.Stats().IdleCount)

	conn := pool.Get()
	fmt.Printf("active:%v, ide: %v\n ", pool.Stats().ActiveCount, pool.Stats().IdleCount)
	_, err := conn.Do("PING")
	fmt.Println("err: ", err)
}
func newPool() {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 10 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}
