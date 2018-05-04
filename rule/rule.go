package rule

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type ruleMap struct {
	sync.RWMutex
	keyA     []string
	keyB     []string
	mapRuleA sync.Map
	mapRuleB sync.Map
	chose    choseInt
}

type choseInt int

const (
	choseAMap choseInt = iota
	choseBMap
)

var rm ruleMap

type Rule struct {
	Id       int32
	Weight   string
	ObjInsCd string
	TranTp   string
	ProdCd   string
	TrnCd    string
	BizCd    string
	MchntGrp string
	MchntCD  string
}

// 匹配复合条件的rule
func GetRule(cond *Rule) (int32, bool) {
	keys, rmap := rm.get()
	for _, key := range keys {
		v, ok := rmap.Load(key)
		if !ok {
			continue
		}
		if isMatch(cond, v.(*Rule)) {
			return v.(*Rule).Id, true
		}
	}
	return 0, false
}

// 加载规则到redis中
func StoreRule(db *gorm.DB) error {
	// 获取暂时不用的缓存
	rows, err := db.Model(&Rule{}).Select("id, weight, obj_ins_cd, tran_tp, prod_cd, trn_cd, biz_cd, mchnt_grp, mchnt_cd").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	keys := make([]string, 0)
	var rmap sync.Map
	for rows.Next() {
		var item Rule
		if err := db.ScanRows(rows, &item); err != nil {
			return err
		}
		keyId := fmt.Sprintf("%v_%v", item.Weight, item.Id)
		keys = append(keys, keyId)
		rmap.Store(keyId, &item)
	}
	sort.Strings(keys)
	rm.set(keys, rmap)
	return nil
}

func (r *ruleMap) set(keys []string, rmap sync.Map) {
	r.Lock()
	defer r.Unlock()

	if r.chose == choseAMap {
		r.chose = choseBMap
		r.keyB = keys
		r.mapRuleB = rmap
	} else {
		r.chose = choseAMap
		r.keyA = keys
		r.mapRuleA = rmap
	}
}

func (r *ruleMap) get() ([]string, sync.Map) {
	r.RLock()
	defer r.RUnlock()

	key := make([]string, 0, len(r.keyA))
	if r.chose == choseAMap {
		key = append(key, r.keyA...)
		return key, r.mapRuleA
	}

	key = append(key, r.keyB...)
	return key, r.mapRuleB
}

func isMatch(cond, val *Rule) bool {
	if !isMemMatch(cond.ObjInsCd, val.ObjInsCd) {
		return false
	}
	if !isMemMatch(cond.TranTp, val.TranTp) {
		return false
	}
	if !isMemMatch(cond.ProdCd, val.ProdCd) {
		return false
	}
	if !isMemMatch(cond.TrnCd, val.TrnCd) {
		return false
	}
	if !isMemMatch(cond.BizCd, val.BizCd) {
		return false
	}
	if !isMemMatch(cond.MchntGrp, val.MchntGrp) {
		return false
	}
	if !isMemMatch(cond.MchntCD, val.MchntCD) {
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
