package bizm

import (
	"net"
	"ninja/base/misc/log"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpc *grpc.Server
}

func (s *GrpcServer) Serve(ln net.Listener) {
	log.Errorf("err: %v", s.GetServer().Serve(ln))
}

func (s *GrpcServer) GetServer() *grpc.Server {
	if s.grpc == nil {
		s.grpc = grpc.NewServer()
	}
	return s.grpc
}
