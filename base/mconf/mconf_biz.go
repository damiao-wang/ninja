package mconf

import (
	"sync"
)

const (
	// DEV is for develop
	DEV = "dev"
	// PROD is for production
	PROD = "prod"
)

type BaseConfig struct {
	Pkg         string            `toml:"pkg"`
	RunMode     string            `toml:"run_mode"`
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

func GetRunMode() string {
	if initBaseConfig().RunMode == "" {
		return DEV
	}
	return initBaseConfig().RunMode
}

func GetPkgName() string {
	if initBaseConfig().Pkg == "" {
		panic("配置文件必须包含pkg字段.")
	}

	return initBaseConfig().Pkg
}
