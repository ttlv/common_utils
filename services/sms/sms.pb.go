// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sms.proto

/*
Package sms is a generated protocol buffer package.

It is generated from these files:
	sms.proto

It has these top-level messages:
	SendParams
	SendResp
*/
package sms

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SendParams struct {
	Brand   string `protobuf:"bytes,1,opt,name=brand" json:"brand,omitempty"`
	Country string `protobuf:"bytes,2,opt,name=country" json:"country,omitempty"`
	Phone   string `protobuf:"bytes,3,opt,name=phone" json:"phone,omitempty"`
	Content string `protobuf:"bytes,4,opt,name=content" json:"content,omitempty"`
}

func (m *SendParams) Reset()                    { *m = SendParams{} }
func (m *SendParams) String() string            { return proto.CompactTextString(m) }
func (*SendParams) ProtoMessage()               {}
func (*SendParams) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SendParams) GetBrand() string {
	if m != nil {
		return m.Brand
	}
	return ""
}

func (m *SendParams) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *SendParams) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *SendParams) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type SendResp struct {
	Uid   string `protobuf:"bytes,1,opt,name=uid" json:"uid,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
}

func (m *SendResp) Reset()                    { *m = SendResp{} }
func (m *SendResp) String() string            { return proto.CompactTextString(m) }
func (*SendResp) ProtoMessage()               {}
func (*SendResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SendResp) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *SendResp) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*SendParams)(nil), "sms.SendParams")
	proto.RegisterType((*SendResp)(nil), "sms.SendResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Sms service

type SmsClient interface {
	Send(ctx context.Context, in *SendParams, opts ...grpc.CallOption) (*SendResp, error)
}

type smsClient struct {
	cc *grpc.ClientConn
}

func NewSmsClient(cc *grpc.ClientConn) SmsClient {
	return &smsClient{cc}
}

func (c *smsClient) Send(ctx context.Context, in *SendParams, opts ...grpc.CallOption) (*SendResp, error) {
	out := new(SendResp)
	err := grpc.Invoke(ctx, "/sms.Sms/Send", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Sms service

type SmsServer interface {
	Send(context.Context, *SendParams) (*SendResp, error)
}

func RegisterSmsServer(s *grpc.Server, srv SmsServer) {
	s.RegisterService(&_Sms_serviceDesc, srv)
}

func _Sms_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SmsServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sms.Sms/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SmsServer).Send(ctx, req.(*SendParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _Sms_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sms.Sms",
	HandlerType: (*SmsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _Sms_Send_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sms.proto",
}

func init() { proto.RegisterFile("sms.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 176 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x8f, 0xc1, 0xaa, 0xc2, 0x30,
	0x10, 0x45, 0xe9, 0x4b, 0x9f, 0xda, 0x01, 0x51, 0x82, 0x8b, 0xe0, 0x4a, 0xba, 0x10, 0x37, 0x76,
	0x51, 0x7f, 0x44, 0xda, 0x2f, 0x68, 0x6d, 0x40, 0xc1, 0x4c, 0xca, 0x4c, 0xba, 0xf0, 0xef, 0x65,
	0x12, 0xab, 0xbb, 0x9c, 0xe1, 0x64, 0xee, 0x5c, 0x28, 0xd8, 0x71, 0x35, 0x92, 0x0f, 0x5e, 0x2b,
	0x76, 0x5c, 0x3e, 0x01, 0x5a, 0x8b, 0xc3, 0xb5, 0xa3, 0xce, 0xb1, 0xde, 0xc1, 0x7f, 0x4f, 0x1d,
	0x0e, 0x26, 0x3b, 0x64, 0xa7, 0xa2, 0x49, 0xa0, 0x0d, 0x2c, 0x6f, 0x7e, 0xc2, 0x40, 0x2f, 0xf3,
	0x17, 0xe7, 0x33, 0x8a, 0x3f, 0xde, 0x3d, 0x5a, 0xa3, 0x92, 0x1f, 0x21, 0xf9, 0x18, 0x2c, 0x06,
	0x93, 0xcf, 0x7e, 0xc4, 0xb2, 0x86, 0x95, 0xa4, 0x35, 0x96, 0x47, 0xbd, 0x05, 0x35, 0x3d, 0xe6,
	0x24, 0x79, 0xca, 0x36, 0x4b, 0xe4, 0xe9, 0x93, 0x92, 0xa0, 0x3e, 0x83, 0x6a, 0x1d, 0xeb, 0x23,
	0xe4, 0xf2, 0x55, 0x6f, 0x2a, 0x69, 0xf0, 0xbb, 0x79, 0xbf, 0xfe, 0x0e, 0x64, 0x6d, 0xbf, 0x88,
	0xe5, 0x2e, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd0, 0x47, 0x2c, 0x6d, 0xe9, 0x00, 0x00, 0x00,
}
