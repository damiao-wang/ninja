package rule

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/garyburd/redigo/redis"
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

func GetRouteCfg(conn redis.Conn, key string, cond *TblRouteCfg) (*TblRouteCfg, error) {
	reply, err := redis.ByteSlices(conn.Do("ZRANGE", key, 0, -1))
	if err != nil {
		return nil, err
	}
	for _, v := range reply {
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
