package mconf

import (
	"io/ioutil"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	defatultConf = "conf/conf.toml"

	flagOnce sync.Once
)

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

func ReadFile(obj interface{}) error {
	data, err := ioutil.ReadFile(defatultConf)
	if err != nil {
		log.Printf("ReadFile, error: %v", err)
		return err
	}
	return Unmarshal(data, obj)
}

func Unmarshal(data []byte, obj interface{}) error {
	if err := toml.Unmarshal(data, obj); err != nil {
		log.Printf("Unmarshal, error: %v", err)
		return err
	}
	return nil
}
