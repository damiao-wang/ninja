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
	WebPort  int
	GrpcPort int
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
	webKey := fmt.Sprintf("service %v", srvName)
	grpcKey := fmt.Sprintf("grpc %v", srvName)

	if v, ok := data[webKey]; ok {
		port, err := c.GetPortInt(v)
		if err != nil {
			return errors.Trace(err)
		}
		c.WebPort = port
	}

	if v, ok := data[grpcKey]; ok {
		port, err := c.GetPortInt(v)
		if err != nil {
			return errors.Trace(err)
		}
		c.GrpcPort = port
	}

	if c.WebPort == 0 || c.GrpcPort == 0 {
		return errors.Fmt(`port of "service:%v" is not found`, srvName)
	}

	return nil
}

func (c *Config) GetPortInt(addr string) (int, error) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, errors.Trace(err)
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return 0, errors.Trace(err)
	}
	return portInt, nil
}

func (c *Config) Close() error {
	if c.logFiler != nil {
		return c.logFiler.Close()
	}
	return nil
}
