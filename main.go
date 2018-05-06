package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

var pool *redis.Pool

func main() {
	// db, _ := gorm.Open("mysql", "root:wang123456@/hyjr?charset=utf8&parseTime=True&loc=Local")
	// err := rule.StoreRoute(db, conn, "rule2")
	// fmt.Println(err)

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

func get() {

}
