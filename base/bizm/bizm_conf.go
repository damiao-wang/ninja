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

func (c *Config) Init(srvName string) error {
	pkgName := mconf.GetPkgName()
	c.RunMode = mconf.GetRunMode()
	if c.RunMode == mconf.PROD {
		log.SetLevel("INFO")
		c.logFiler = &lumberjack.Logger{
			Filename:   fmt.Sprintf("/var/logs/%v/%v.log", pkgName, srvName),
			MaxAge:     30,
			MaxBackups: 5,
		}
		log.SetOut(c.logFiler)
	}
	return c.setPortByName(srvName)
}

func (c *Config) setPortByName(srvName string) error {
	data := mconf.GetListenAddrs()
	srvKey := fmt.Sprintf("service %v", srvName)

	if v, ok := data[srvKey]; ok {
		return c.SetPortInt(v)
	}

	return errors.Fmt(`port of "service %v" is not found`, srvName)
}

func (c *Config) SetPortInt(addr string) error {
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
	if c.logFiler != nil {
		return c.logFiler.Close()
	}
	return nil
}
