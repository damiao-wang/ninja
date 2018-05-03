package bizm

import (
	"fmt"
	"net"

	"ninja/base/misc/context"
	"ninja/base/misc/errors"

	"github.com/soheilhy/cmux"
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

func (s *Server) Run(ctx context.T) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Conf.Port))
	if err != nil {
		errors.Trace(err)
	}

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

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

	return handleLnError(m.Serve())
}

func (s *Server) Close() error {
	err := s.WebServer.Close()
	err = s.Conf.Close()
	return err
}

func handleLnError(err error) error {
	if err, ok := err.(*net.OpError); ok {
		if err.Op != "accept" {
			return nil
		}
		if err.Err.Error() == "use of closed network connection" {
			return nil
		}
	}
	return err
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
