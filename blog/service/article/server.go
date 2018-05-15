package article

import (
	"bytes"
	"fmt"
	"ninja/base/bizm"
	"ninja/base/misc/context"
	"ninja/base/misc/errors"
	pb "ninja/blog/rpc/blog"

	"github.com/tealeg/xlsx"
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

func (s *Service) Hello(ctx context.T, req *pb.HelloReq) (*pb.HelloResp, error) {
	msg := fmt.Sprintf("Hello %v.", req.Name)
	return &pb.HelloResp{
		Msg: msg,
	}, nil
}

func (s *Service) ExportInfo(ctx context.T, req *pb.ExportInfoReq) (*pb.ExportInfoResp, error) {
	w := xlsx.NewFile()
	sheet, err := w.AddSheet("sheet")
	if err != nil {
		return nil, errors.Trace(err)
	}
	row := sheet.AddRow()
	row.AddCell().SetString("ID")
	row.AddCell().SetString("NAME")
	row = sheet.AddRow()
	row.AddCell().SetInt(1)
	row.AddCell().SetString(req.Name)

	buf := bytes.NewBuffer(nil)
	w.Write(buf)
	return &pb.ExportInfoResp{
		Filename: "Info.xlsx",
		Data:     buf.Bytes(),
	}, nil
}
