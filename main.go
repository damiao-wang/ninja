package main

import (
	"fmt"
	"ninja/rule"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// db, _ := gorm.Open("mysql", "root:wang123456@/hyjr?charset=utf8&parseTime=True&loc=Local")
	conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
	// err := rule.StoreRoute(db, conn, "rule2")
	// fmt.Println(err)

	r, err := rule.GetRouteCfg(conn, "rule2", &rule.Rule{ProdCd: "12345", MchntCD: "1"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
}
