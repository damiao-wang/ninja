package bizm

import (
	"net"
	"strconv"

	"ninja/base/misc/errors"

	"golang.org/x/net/context"
)

type Server struct {
	Port int
	WebServer
}

func (s *Server) Init(srv interface{}, pkg string, register func() error) error {
	if err := SetPortByName(pkg, s); err != nil {
		return errors.Trace(err)
	}
	if err := register(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (s *Server) SetPort(addr string) error {
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

func (s *Server) Run(ctx context.Context) error {
	return s.WebServer.Serve(s.Port)
}

func (s *Server) Close() error {
	return nil
}
