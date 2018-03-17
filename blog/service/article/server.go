package article

import (
	"fmt"

	"ninja/base/bizm"
	"golang.org/x/net/context"
)

type Service struct{
	bizm.Service
}

// func (s *Service) Bye(ctx *gin.Context) {
// 	name, _ := ctx.Get("name")
// 	ctx.String(http.StatusOK, "Bye %v.", name)
// }

func (s *Service) Desc() string {
	return "Article"
}

func (s *Service) Register() error {
	
	return nil
}

type HelloReq struct {
	Name string
}

type HelloResp struct {
	Msg string
}

func (s *Service) Hello(c context.Context, req HelloReq) (*HelloResp, error) {
	msg := fmt.Sprintf("Hello %v.", req.Name)
	return &HelloResp{
		Msg: msg,
	}, nil
}
