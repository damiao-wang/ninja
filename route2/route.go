package route2

import (
	"errors"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
)

const (
	routeKey  = "tbl_route_cfg_key"
	routeKeyA = "tbl_route_cfg_key_a"
	routeKeyB = "tbl_route_cfg_key_b"
)

var redisPool *redis.Pool

// 匹配复合条件的rule
func GetRoute(cond *TblRouteCfg) (*TblRouteCfg, error) {
	routes, err := getRoutes()
	if err != nil {
		return nil, err
	}
	for _, buf := range routes {
		item := &TblRouteCfg{}
		if err := proto.Unmarshal(buf, item); err != nil {
			return nil, err
		}
		if isMatch(cond, item) {
			return item, nil
		}
	}
	return nil, errors.New("Find nothing!")
}

// 加载配置到redis中
func LoadRoute(db *gorm.DB) error {
	// 获取暂时不用的缓存
	var routes []*TblRouteCfg
	if err := db.Table("tbl_route_cfg").Find(&routes).Error; err != nil {
		return err
	}

	return setRoutes(routes)
}

func InitRedis(addr string) {
	redisPool = &redis.Pool{
		MaxIdle:     8,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
}

func setRoutes(routes []*TblRouteCfg) error {
	conn := redisPool.Get()
	defer conn.Close()

	key, err := redis.String(conn.Do("GET", routeKey))
	if err != nil {
		return err
	}

	var anotherKey string
	if routeKeyA == key {
		anotherKey = routeKeyB
	} else {
		anotherKey = routeKeyA
	}

	conn.Send("DEL", anotherKey, anotherKey)
	for _, v := range routes {
		buf, err := proto.Marshal(v)
		if err != nil {
			return err
		}
		conn.Send("ZADD", anotherKey, v.Proi, buf)
	}
	conn.Send("SET", routeKey, anotherKey)
	conn.Flush()
	_, err = conn.Receive()
	return err

}

func getRoutes() ([][]byte, error) {
	conn := redisPool.Get()
	defer conn.Close()

	key, err := redis.String(conn.Do("GET", routeKey))
	if err != nil {
		return nil, err
	}
	return redis.ByteSlices(conn.Do("ZRANGE", key, 0, -1))
}

func isMatch(cond, val *TblRouteCfg) bool {
	if !isMemMatch(cond.ProdCd, val.ProdCd) {
		return false
	}
	if !isMemMatch(cond.TranCd, val.TranCd) {
		return false
	}
	if !isMemMatch(cond.AppId, val.AppId) {
		return false
	}
	if !isMemMatch(cond.Status, val.Status) {
		return false
	}
	if !isMemMatch(cond.IssInsGrp, val.IssInsGrp) {
		return false
	}
	if !isMemMatch(cond.IssInsCd, val.IssInsCd) {
		return false
	}
	if !isMemMatch(cond.CardBin, val.CardBin) {
		return false
	}
	if !isMemMatch(cond.CardClass, val.CardClass) {
		return false
	}
	if !isMemMatch(cond.AmtL, val.AmtL) {
		return false
	}
	if !isMemMatch(cond.AmtH, val.AmtH) {
		return false
	}
	if !isMemMatch(cond.DateB, val.DateB) {
		return false
	}
	if !isMemMatch(cond.DateE, val.DateE) {
		return false
	}
	if !isMemMatch(cond.TimeB, val.TimeB) {
		return false
	}
	if !isMemMatch(cond.TimeE, val.TimeE) {
		return false
	}
	if !isMemMatch(cond.ObjServer, val.ObjServer) {
		return false
	}
	if !isMemMatch(cond.ObjMchntCd, val.ObjMchntCd) {
		return false
	}
	if !isMemMatch(cond.Use, val.Use) {
		return false
	}
	return true
}

func isMemMatch(src, prefix string) bool {
	if prefix != "*" {
		if !strings.HasPrefix(src, prefix) {
			return false
		}
	}
	return true
}
