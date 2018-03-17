package mconf

import (
	"sync"
)

type BaseConfig struct {
	ListenAddrs map[string]string `toml:"listen_addrs"`
}

var baseConfigInstance *BaseConfig
var baseConfigOnce sync.Once

func initBaseConfig() *BaseConfig {
	baseConfigOnce.Do(func() {
		ByFlag(&baseConfigInstance)
	})
	return baseConfigInstance
}

func GetListenAddrs() map[string]string {
	return initBaseConfig().ListenAddrs
}
