package bizm

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"

	"ninja/base/misc/errors"

	"golang.org/x/net/context"
)

type Server struct {
	Conf *Config
	GrpcServer
	WebServer
}

func (s *Server) Init(srv interface{}, pkg string, register func() error) error {
	s.Conf = &Config{}
	if err := s.Conf.Init(pkg); err != nil {
		return errors.Trace(err)
	}
	if err := register(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (s *Server) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Conf.Port))
	if err != nil {
		return err
	}
	go s.GrpcServer.Serve(ln)
	go s.WebServer.Serve(ln)
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, os.Kill)
	<-done
	return nil
}

func (s *Server) Close() error {
	return s.Conf.Close()
}

func (s *Server) RegisterServer(controller, grpcRegister interface{}) {
	s.AutoRouter(controller)
	fnVal := reflect.ValueOf(grpcRegister)
	fnType := fnVal.Type()
	if fnType.Kind() != reflect.Func {
		panic("grpcRegister must be a func.")
	}
	grpcServer := s.GetServer()
	fnVal.Call([]reflect.Value{
		reflect.ValueOf(grpcServer),
		reflect.ValueOf(controller),
	})
}
