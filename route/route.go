package route

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type routeInfo struct {
	sync.RWMutex
	routeA []*TblRouteCfg
	routeB []*TblRouteCfg
}

const (
	routeKey  = "tbl_route_cfg_key"
	routeKeyA = "tbl_route_cfg_key_a"
	routeKeyB = "tbl_route_cfg_key_b"
)

var (
	rInfo     routeInfo
	redisPool *redis.Pool
)

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

type routeSlice []*TblRouteCfg

func (r routeSlice) Len() int {
	return len(r)
}

func (r routeSlice) Less(i, j int) bool {
	if r[i].Proi == r[j].Proi {
		fmt.Println("In..")
		return r[i].Id < r[j].Id
	}

	return r[i].Proi < r[j].Proi
}

func (r routeSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// 匹配复合条件的rule
func GetRule(cond *TblRouteCfg) (int32, error) {
	routes, err := rInfo.get()
	if err != nil {
		return 0, err
	}
	for _, v := range routes {
		if isMatch(cond, v) {
			return v.Id, nil
		}
	}
	return 0, errors.New("Find nothing!")
}

// 加载配置到redis中
func LoadRule(db *gorm.DB) error {
	// 获取暂时不用的缓存
	var routes []*TblRouteCfg
	if err := db.Find(&routes).Error; err != nil {
		return err
	}

	sort.Sort(routeSlice(routes))
	return rInfo.set(routes)
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

func (r *routeInfo) set(routes []*TblRouteCfg) error {
	key, err := getRouteKey()
	if err != nil {
		return err
	}
	var anotherKey string
	r.Lock()
	if routeKeyA == key {
		anotherKey = routeKeyB
		r.routeB = routes
	} else {
		anotherKey = routeKeyA
		r.routeA = routes
	}
	r.Unlock()

	return setRouteKey(anotherKey)
}

func (r *routeInfo) get() ([]*TblRouteCfg, error) {
	key, err := getRouteKey()
	if err != nil {
		return nil, err
	}
	routes := make([]*TblRouteCfg, 0, len(r.routeA))
	r.RLock()
	if routeKeyA == key {
		routes = append(routes, r.routeA...)
	} else {
		routes = append(routes, r.routeB...)
	}
	r.RUnlock()

	return routes, nil
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

func getRouteKey() (string, error) {
	conn := redisPool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", routeKey))
}

func setRouteKey(val string) error {
	conn := redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", routeKey, val)
	return err
}
