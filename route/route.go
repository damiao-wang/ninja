package route

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	routeKey  = "tbl_route_cfg_key"
	routeKeyA = "tbl_route_cfg_key_a"
	routeKeyB = "tbl_route_cfg_key_b"
)

var redisPool *redis.Pool

type TblRouteCfg struct {
	Id         int32
	Proi       int32
	ProdCd     string `gorm:"default:'*'"`
	TranCd     string `gorm:"default:'*'"`
	AppId      string `gorm:"default:'*'"`
	Status     string `gorm:"default:'*'"`
	IssInsGrp  string `gorm:"default:'*'"`
	IssInsCd   string `gorm:"default:'*'"`
	CardBin    string `gorm:"default:'*'"`
	CardClass  string `gorm:"default:'*'"`
	AmtL       string `gorm:"default:'*'"`
	AmtH       string `gorm:"default:'*'"`
	DateB      string `gorm:"default:'*'"`
	DateE      string `gorm:"default:'*'"`
	TimeB      string `gorm:"default:'*'"`
	TimeE      string `gorm:"default:'*'"`
	ObjServer  string `gorm:"default:'*'"`
	ObjMchntCd string
	Use        string `gorm:"default:'1'"`
}

// 匹配复合条件的rule
func GetRoute(cond *TblRouteCfg) (*TblRouteCfg, error) {
	routes, err := getRoutes()
	if err != nil {
		return nil, err
	}
	for _, v := range routes {
		var item *TblRouteCfg
		if err := json.Unmarshal(v, &item); err != nil {
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
	if err := db.Find(&routes).Error; err != nil {
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
		item, err := json.Marshal(v)
		if err != nil {
			return err
		}
		conn.Send("ZADD", anotherKey, v.Proi, item)
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

func (TblRouteCfg) TableName() string {
	return "tbl_route_cfg"
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
