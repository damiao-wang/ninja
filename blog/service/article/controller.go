package article

import (
	"fmt"

	"golang.org/x/net/context"
)

type Controller struct{}

type HelloReq struct {
	Name string
}

type HelloResp struct {
	Msg string
}

func (c *Controller) Hello(ctx context.Context, req HelloReq) (*HelloResp, error) {
	msg := fmt.Sprintf("Hello %v.", req.Name)
	return &HelloResp{
		Msg: msg,
	}, nil
}
