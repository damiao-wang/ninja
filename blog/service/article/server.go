package article

import (
	"ninja/base/bizm"
	pb "ninja/blog/rpc/blog"
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
