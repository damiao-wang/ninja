package bizm

import (
	"net"
	"strconv"
	"gitlab.1dmy.com/ezbuy/base/misc/errors"
)

type Service struct {
	Port int
}

func (s *Service) Init(srv interface{}, pkg string, register func() error) error {
	if err := SetPortByName(pkg, s); err != nil {
		return errors.Trace(err)
	}
	if err := register(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (s *Service) SetPort(addr string) error {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return errors.Trace(err)
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return errors.Trace(err)
	}
	s.Port = portInt
	return nil
}

func (s *Service) Close() error {
	return nil
}
