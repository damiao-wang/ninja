package bizm

import (
	"net"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpc *grpc.Server
}

func (s *GrpcServer) Serve(ln net.Listener) {
	s.GetServer().Serve(ln)
}

func (s *GrpcServer) GetServer() *grpc.Server {
	if s.grpc == nil {
		s.grpc = grpc.NewServer()
	}
	return s.grpc
}
