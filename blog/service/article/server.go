package article

import (
	"bytes"
	"fmt"
	"io"
	"ninja/base/bizm"
	"ninja/base/misc/context"
	"ninja/base/misc/errors"
	pb "ninja/blog/rpc/blog"
	"os"

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

func (s *Service) Upload(ctx context.T, req *pb.UploadReq) (*pb.UploadResp, error) {
	file, err := os.Create(req.Filename)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer file.Close()

	buf := bytes.NewBuffer(req.Data)
	if _, err := io.Copy(file, buf); err != nil {
		return nil, errors.Trace(err)
	}

	return &pb.UploadResp{}, nil
}
