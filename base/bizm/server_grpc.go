package bizm

import (
	"net"
	"ninja/base/misc/log"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpc *grpc.Server
}

func (s *GrpcServer) Serve(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s.GetServer().Serve(ln)
}

func (s *GrpcServer) GetServer() *grpc.Server {
	if s.grpc == nil {
		s.grpc = grpc.NewServer()
	}
	return s.grpc
}
