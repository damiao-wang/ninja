package article

import (
	"fmt"

	pb "ninja/blog/rpc/blog"

	"golang.org/x/net/context"
)

type Controller struct{}

//@router [get]
func (c *Controller) Hello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	msg := fmt.Sprintf("Hello %v.", req.Name)
	return &pb.HelloResp{
		Msg: msg,
	}, nil
}
