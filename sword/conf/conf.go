package conf

import (
	"time"

	"ninja/base/mconf"
)

type Config struct {
	Title    string
	Hosts    []string
	Owner    Owner
	Database Database
	Servers  Servers
	Books    []string
	Clients  Clients
	Products []Product
}

type Owner struct {
	Name string
	Dob  time.Time
}

type Database struct {
	Server        string
	Ports         []int
	ConnectionMax int `toml:"connection_max"`
	Enabled       bool
}

type Servers struct {
	Alpha Sinfo
	Beta  Sinfo
}

type Sinfo struct {
	Ip string
	Dc string
}

type Clients struct {
	Data [][]interface{}
}

type Product struct {
	Name string
	Sku  int
}

var instance Config

func Get() *Config {
	mconf.ByFlagOnce(&instance)
	return &instance
}
