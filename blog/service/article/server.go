package article

import (
	"ninja/base/bizm"
)

type Service struct{
	bizm.Server
}

func (s *Service) Desc() string {
	return "Article"
}

func (s *Service) Register() error {
	s.AutoRouter(&Controller{})
	return nil
}
