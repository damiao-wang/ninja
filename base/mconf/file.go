package mconf

import (
	"io/ioutil"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/spf13/pflag"
)

var (
	defatultConf = "conf/conf.toml"

	flagOnce sync.Once
)

func ReadFile(obj interface{}) error {
	data, err := ioutil.ReadFile(defatultConf)
	if err != nil {
		log.Printf("ReadFile, error: %v", err)
		return err
	}
	return Unmarshal(data, obj)
}

func ByFlagOnce(obj interface{}) {
	flagOnce.Do(func() {
		err := ReadFile(obj)
		if err != nil {
			log.Printf("ByFlagOnce, error: %v", err)
			panic(err)
		}
	})
	return
}

func ByFlag(obj interface{}) {
	err := ReadFile(obj)
	if err != nil {
		panic(err)
	}
}

func Unmarshal(data []byte, obj interface{}) error {
	if err := toml.Unmarshal(data, obj); err != nil {
		log.Printf("Unmarshal, error: %v", err)
		return err
	}
	return nil
}

func ParseFlag(p *pflag.FlagSet) {
	hasConf := false
	p.VisitAll(func(f *pflag.Flag) {
		if f.Name == "conf" {
			hasConf = true
		}
		if f.Shorthand == "c" {
			hasConf = true
		}
	})
	if !hasConf {
		p.StringVarP(&defatultConf, "conf", "c", defatultConf, "ctmpl path")
	}
}
