package bizm

import (
	"fmt"
	"net"

	"ninja/base/misc/errors"

	"github.com/soheilhy/cmux"
	"golang.org/x/net/context"
)

type Server struct {
	Conf *Config
	GrpcServer
	WebServer
}

func (s *Server) Init(srv interface{}, srvName string, register func() error) error {
	s.Conf = &Config{}
	if err := s.Conf.Init(srvName); err != nil {
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
		errors.Trace(err)
	}

	m := cmux.New(ln)
	// start grpc
	{
		match := cmux.HTTP2HeaderField("content-type", "application/grpc")
		ln := m.Match(match)
		go s.GrpcServer.Serve(ln)

	}

	// start webapi
	{
		ln := m.Match(cmux.Any())
		go s.WebServer.Serve(ln)
	}

	return m.Serve()
}

func (s *Server) Close() error {
	return s.Conf.Close()
}

// func (s *Server) RegisterServer(controller, grpcRegister interface{}) {
// 	s.AutoRouter(controller)
// 	fnVal := reflect.ValueOf(grpcRegister)
// 	fnType := fnVal.Type()
// 	if fnType.Kind() != reflect.Func {
// 		panic("grpcRegister must be a func.")
// 	}
// 	grpcServer := s.GetServer()
// 	fnVal.Call([]reflect.Value{
// 		reflect.ValueOf(grpcServer),
// 		reflect.ValueOf(controller),
// 	})
// }
