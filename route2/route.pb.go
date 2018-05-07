// Code generated by protoc-gen-go. DO NOT EDIT.
// source: route.proto

/*
Package route2 is a generated protocol buffer package.

It is generated from these files:
	route.proto

It has these top-level messages:
	TblRouteCfg
*/
package route2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TblRouteCfg struct {
	Id         int32  `protobuf:"varint,1,opt,name=Id" json:"Id,omitempty"`
	Proi       int32  `protobuf:"varint,2,opt,name=Proi" json:"Proi,omitempty"`
	ProdCd     string `protobuf:"bytes,3,opt,name=ProdCd" json:"ProdCd,omitempty"`
	TranCd     string `protobuf:"bytes,4,opt,name=TranCd" json:"TranCd,omitempty"`
	AppId      string `protobuf:"bytes,5,opt,name=AppId" json:"AppId,omitempty"`
	Status     string `protobuf:"bytes,6,opt,name=Status" json:"Status,omitempty"`
	IssInsGrp  string `protobuf:"bytes,7,opt,name=IssInsGrp" json:"IssInsGrp,omitempty"`
	IssInsCd   string `protobuf:"bytes,8,opt,name=IssInsCd" json:"IssInsCd,omitempty"`
	CardBin    string `protobuf:"bytes,9,opt,name=CardBin" json:"CardBin,omitempty"`
	CardClass  string `protobuf:"bytes,10,opt,name=CardClass" json:"CardClass,omitempty"`
	AmtL       string `protobuf:"bytes,11,opt,name=AmtL" json:"AmtL,omitempty"`
	AmtH       string `protobuf:"bytes,12,opt,name=AmtH" json:"AmtH,omitempty"`
	DateB      string `protobuf:"bytes,13,opt,name=DateB" json:"DateB,omitempty"`
	DateE      string `protobuf:"bytes,14,opt,name=DateE" json:"DateE,omitempty"`
	TimeB      string `protobuf:"bytes,15,opt,name=TimeB" json:"TimeB,omitempty"`
	TimeE      string `protobuf:"bytes,16,opt,name=TimeE" json:"TimeE,omitempty"`
	ObjServer  string `protobuf:"bytes,17,opt,name=ObjServer" json:"ObjServer,omitempty"`
	ObjMchntCd string `protobuf:"bytes,18,opt,name=ObjMchntCd" json:"ObjMchntCd,omitempty"`
	Use        string `protobuf:"bytes,19,opt,name=Use" json:"Use,omitempty"`
}

func (m *TblRouteCfg) Reset()                    { *m = TblRouteCfg{} }
func (m *TblRouteCfg) String() string            { return proto.CompactTextString(m) }
func (*TblRouteCfg) ProtoMessage()               {}
func (*TblRouteCfg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TblRouteCfg) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *TblRouteCfg) GetProi() int32 {
	if m != nil {
		return m.Proi
	}
	return 0
}

func (m *TblRouteCfg) GetProdCd() string {
	if m != nil {
		return m.ProdCd
	}
	return ""
}

func (m *TblRouteCfg) GetTranCd() string {
	if m != nil {
		return m.TranCd
	}
	return ""
}

func (m *TblRouteCfg) GetAppId() string {
	if m != nil {
		return m.AppId
	}
	return ""
}

func (m *TblRouteCfg) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *TblRouteCfg) GetIssInsGrp() string {
	if m != nil {
		return m.IssInsGrp
	}
	return ""
}

func (m *TblRouteCfg) GetIssInsCd() string {
	if m != nil {
		return m.IssInsCd
	}
	return ""
}

func (m *TblRouteCfg) GetCardBin() string {
	if m != nil {
		return m.CardBin
	}
	return ""
}

func (m *TblRouteCfg) GetCardClass() string {
	if m != nil {
		return m.CardClass
	}
	return ""
}

func (m *TblRouteCfg) GetAmtL() string {
	if m != nil {
		return m.AmtL
	}
	return ""
}

func (m *TblRouteCfg) GetAmtH() string {
	if m != nil {
		return m.AmtH
	}
	return ""
}

func (m *TblRouteCfg) GetDateB() string {
	if m != nil {
		return m.DateB
	}
	return ""
}

func (m *TblRouteCfg) GetDateE() string {
	if m != nil {
		return m.DateE
	}
	return ""
}

func (m *TblRouteCfg) GetTimeB() string {
	if m != nil {
		return m.TimeB
	}
	return ""
}

func (m *TblRouteCfg) GetTimeE() string {
	if m != nil {
		return m.TimeE
	}
	return ""
}

func (m *TblRouteCfg) GetObjServer() string {
	if m != nil {
		return m.ObjServer
	}
	return ""
}

func (m *TblRouteCfg) GetObjMchntCd() string {
	if m != nil {
		return m.ObjMchntCd
	}
	return ""
}

func (m *TblRouteCfg) GetUse() string {
	if m != nil {
		return m.Use
	}
	return ""
}

func init() {
	proto.RegisterType((*TblRouteCfg)(nil), "route2.TblRouteCfg")
}

func init() { proto.RegisterFile("route.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 300 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x91, 0xdd, 0x4e, 0xf2, 0x40,
	0x10, 0x86, 0xc3, 0x5f, 0x81, 0xe1, 0xfb, 0x10, 0x47, 0x63, 0x26, 0xc6, 0x18, 0xe2, 0x11, 0x47,
	0x1e, 0xe8, 0x15, 0xc0, 0xda, 0x48, 0x13, 0x0d, 0x04, 0xea, 0x05, 0xb4, 0xee, 0xaa, 0x25, 0xd0,
	0x36, 0xbb, 0x8b, 0xf7, 0xe9, 0x1d, 0x99, 0x9d, 0xa5, 0xad, 0x67, 0xef, 0xf3, 0xbc, 0x9b, 0xc9,
	0x4c, 0x16, 0x46, 0xba, 0x38, 0x5a, 0x75, 0x5f, 0xea, 0xc2, 0x16, 0x18, 0x30, 0x3c, 0xdc, 0xfd,
	0x74, 0x60, 0x14, 0xa7, 0xfb, 0x8d, 0x23, 0xf1, 0xf1, 0x89, 0x63, 0x68, 0x47, 0x92, 0x5a, 0xd3,
	0xd6, 0xac, 0xb7, 0x69, 0x47, 0x12, 0x11, 0xba, 0x6b, 0x5d, 0x64, 0xd4, 0x66, 0xc3, 0x19, 0xaf,
	0x20, 0x58, 0xeb, 0x42, 0x0a, 0x49, 0x9d, 0x69, 0x6b, 0x36, 0xdc, 0x9c, 0xc8, 0xf9, 0x58, 0x27,
	0xb9, 0x90, 0xd4, 0xf5, 0xde, 0x13, 0x5e, 0x42, 0x6f, 0x5e, 0x96, 0x91, 0xa4, 0x1e, 0x6b, 0x0f,
	0xee, 0xf5, 0xd6, 0x26, 0xf6, 0x68, 0x28, 0xf0, 0xaf, 0x3d, 0xe1, 0x0d, 0x0c, 0x23, 0x63, 0xa2,
	0xdc, 0x3c, 0xeb, 0x92, 0xfa, 0x5c, 0x35, 0x02, 0xaf, 0x61, 0xe0, 0x41, 0x48, 0x1a, 0x70, 0x59,
	0x33, 0x12, 0xf4, 0x45, 0xa2, 0xe5, 0x22, 0xcb, 0x69, 0xc8, 0x55, 0x85, 0x6e, 0xa6, 0x8b, 0x62,
	0x9f, 0x18, 0x43, 0xe0, 0x67, 0xd6, 0xc2, 0xdd, 0x38, 0x3f, 0xd8, 0x17, 0x1a, 0x71, 0xc1, 0xf9,
	0xe4, 0x96, 0xf4, 0xaf, 0x76, 0x4b, 0x77, 0xc7, 0x53, 0x62, 0xd5, 0x82, 0xfe, 0xfb, 0x3b, 0x18,
	0x2a, 0x1b, 0xd2, 0xb8, 0xb1, 0xa1, 0xb3, 0x71, 0x76, 0x50, 0x0b, 0x3a, 0xf3, 0x96, 0xa1, 0xb2,
	0x21, 0x4d, 0x1a, 0x1b, 0xba, 0xed, 0x56, 0xe9, 0x6e, 0xab, 0xf4, 0xb7, 0xd2, 0x74, 0xee, 0xb7,
	0xab, 0x05, 0xde, 0x02, 0xac, 0xd2, 0xdd, 0xeb, 0xfb, 0x57, 0x6e, 0x85, 0x24, 0xe4, 0xfa, 0x8f,
	0xc1, 0x09, 0x74, 0xde, 0x8c, 0xa2, 0x0b, 0x2e, 0x5c, 0x4c, 0x03, 0xfe, 0xe2, 0xc7, 0xdf, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x44, 0x4f, 0xab, 0x72, 0xf1, 0x01, 0x00, 0x00,
}
