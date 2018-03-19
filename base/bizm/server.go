package bizm

import (
	"ninja/base/misc/errors"

	"golang.org/x/net/context"
)

type Server struct {
	Conf *Config
	WebServer
}

func (s *Server) Init(srv interface{}, pkg string, register func() error) error {
	if err := s.Conf.Init(pkg); err != nil {
		return errors.Trace(err)
	}
	if err := register(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (s *Server) Run(ctx context.Context) error {
	return s.WebServer.Serve(s.Conf.Port)
}

func (s *Server) Close() error {
	return s.Conf.Close()
}
