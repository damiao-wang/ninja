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
	chose  choseInt
	ruleAs []*Rule
	ruleBs []*Rule
}

type choseInt int

const (
	choseA choseInt = iota
	choseB
)

var rm ruleMap

type Rule struct {
	Id       int32
	Weight   int32
	ObjInsCd string
	TranTp   string
	ProdCd   string
	TrnCd    string
	BizCd    string
	MchntGrp string
	MchntCD  string
}

type ruleSlice []*Rule

func (r ruleSlice) Len() int {
	return len(r)
}

func (r ruleSlice) Less(i, j int) bool {
	if r[i].Weight == r[j].Weight {
		fmt.Println("In..")
		return r[i].Id < r[j].Id
	}

	return r[i].Weight < r[j].Weight
}

func (r ruleSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// 匹配复合条件的rule
func GetRule(cond *Rule) (int32, bool) {
	rules := rm.get()
	for _, v := range rules {
		if isMatch(cond, v) {
			return v.Id, true
		}
	}
	return 0, false
}

// 加载规则到redis中
func StoreRule(db *gorm.DB) error {
	// 获取暂时不用的缓存
	var rules []*Rule
	if err := db.Find(&rules).Error; err != nil {
		return err
	}

	sort.Sort(ruleSlice(rules))
	rm.set(rules)
	return nil
}

func (r *ruleMap) set(rules []*Rule) {
	r.Lock()
	defer r.Unlock()

	if r.chose == choseA {
		r.chose = choseB
		r.ruleBs = rules
	} else {
		r.chose = choseA
		r.ruleAs = rules
	}
}

func (r *ruleMap) get() []*Rule {
	r.RLock()
	defer r.RUnlock()

	rules := make([]*Rule, 0, len(r.ruleAs))
	if r.chose == choseA {
		rules = append(rules, r.ruleAs...)
	} else {
		rules = append(rules, r.ruleBs...)
	}

	return rules
}

func (Rule) TableName() string {
	return "rule"
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
