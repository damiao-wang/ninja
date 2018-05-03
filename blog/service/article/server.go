package article

import (
	"context"
	"errors"
	"fmt"
	"ninja/base/bizm"
	nerrors "ninja/base/misc/errors"
	pb "ninja/blog/rpc/blog"
)

type Service struct {
	bizm.Server
}

func (s *Service) Desc() string {
	return "Article"
}

func (s *Service) Register() error {
	s.InitRouter()
	// pb.RegisterArticleServer(s.GetServer(), s)
	return nil
}

func (s *Service) Hello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	err := PPP()
	return nil, nerrors.Trace(err)

	msg := fmt.Sprintf("Hello %v.", req.Name)
	return &pb.HelloResp{
		Msg: msg,
	}, nil
}

// func init() {
// 	raven.SetDSN("https://1c115aaf0b2048f485936409b03ce0f7:c6480693facf40d6b89c28711fb363b9@sentry.io/304312")
// }

func PPP() error {
	return nerrors.Trace(errors.New("abc"))
}
