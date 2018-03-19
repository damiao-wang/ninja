package bizm

import (
	"fmt"
	"io"
	"net"
	"strconv"

	"ninja/base/mconf"
	"ninja/base/misc/errors"
	"ninja/base/misc/log"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Port     int
	RunMode  string
	logFiler io.WriteCloser
}

func (c *Config) Init(pkgName string) error {
	c.RunMode = mconf.GetRunMode()
	if c.RunMode == mconf.PROD {
		log.SetLevel("INFO")
		c.logFiler = &lumberjack.Logger{
			Filename:   fmt.Sprintf("/var/logs/%v.log", pkgName),
			MaxAge:     30,
			MaxBackups: 5,
		}
		log.SetOut(c.logFiler)
	}
	return c.setPortByName(pkgName)
}

func (c *Config) setPortByName(pkgName string) error {
	data := mconf.GetListenAddrs()
	serviceKey := fmt.Sprintf("service %v", pkgName)
	if pkgName == "main" {
		serviceKey = ""
	}
	if v, ok := data[serviceKey]; ok {
		return c.SetPort(v)
	}
	return errors.Fmt(`port of "service:%v" is not found`, pkgName)

}

func (c *Config) SetPort(addr string) error {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return errors.Trace(err)
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return errors.Trace(err)
	}
	c.Port = portInt
	return nil
}

func (c *Config) Close() error {
	return c.logFiler.Close()
}
