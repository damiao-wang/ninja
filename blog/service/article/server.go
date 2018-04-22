package article

import (
	"fmt"
	"ninja/base/bizm"
	pb "ninja/blog/rpc/blog"

	"golang.org/x/net/context"
)

type Service struct {
	bizm.Server
}

func (s *Service) Desc() string {
	return "Article"
}

func (s *Service) Register() error {
	s.RegisterServer(&Controller{}, pb.RegisterArticleServer)
	return nil
}

//@router [get]
func (s *Service) Hello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	msg := fmt.Sprintf("Hello %v.", req.Name)
	return &pb.HelloResp{
		Msg: msg,
	}, nil
}
