// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/gov/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryEGFParamsRequest is the request type for the Query/EGFParams RPC method.
type QuerySwitchParamsRequest struct {
}

func (m *QuerySwitchParamsRequest) Reset()         { *m = QuerySwitchParamsRequest{} }
func (m *QuerySwitchParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QuerySwitchParamsRequest) ProtoMessage()    {}
func (*QuerySwitchParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_47cb083fb0607b21, []int{0}
}
func (m *QuerySwitchParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuerySwitchParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuerySwitchParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuerySwitchParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuerySwitchParamsRequest.Merge(m, src)
}
func (m *QuerySwitchParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QuerySwitchParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QuerySwitchParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QuerySwitchParamsRequest proto.InternalMessageInfo

// QueryEGFParamsResponse is the response type for the Query/EGFParams RPC
// method.
type QuerySwitchParamsResponse struct {
	Params SwitchParams `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QuerySwitchParamsResponse) Reset()         { *m = QuerySwitchParamsResponse{} }
func (m *QuerySwitchParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QuerySwitchParamsResponse) ProtoMessage()    {}
func (*QuerySwitchParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_47cb083fb0607b21, []int{1}
}
func (m *QuerySwitchParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuerySwitchParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuerySwitchParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuerySwitchParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuerySwitchParamsResponse.Merge(m, src)
}
func (m *QuerySwitchParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QuerySwitchParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QuerySwitchParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QuerySwitchParamsResponse proto.InternalMessageInfo

func (m *QuerySwitchParamsResponse) GetParams() SwitchParams {
	if m != nil {
		return m.Params
	}
	return SwitchParams{}
}

// QueryCustomParamsRequest is the request type for the Query/CustomParams RPC
// method.
type QueryCustomParamsRequest struct {
	MsgUrl string `protobuf:"bytes,1,opt,name=msg_url,json=msgUrl,proto3" json:"msg_url,omitempty"`
}

func (m *QueryCustomParamsRequest) Reset()         { *m = QueryCustomParamsRequest{} }
func (m *QueryCustomParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryCustomParamsRequest) ProtoMessage()    {}
func (*QueryCustomParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_47cb083fb0607b21, []int{2}
}
func (m *QueryCustomParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCustomParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCustomParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCustomParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCustomParamsRequest.Merge(m, src)
}
func (m *QueryCustomParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryCustomParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCustomParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCustomParamsRequest proto.InternalMessageInfo

func (m *QueryCustomParamsRequest) GetMsgUrl() string {
	if m != nil {
		return m.MsgUrl
	}
	return ""
}

// QueryCustomParamsResponse is the response type for the Query/CustomParams RPC
// method.
type QueryCustomParamsResponse struct {
	Params CustomParams `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryCustomParamsResponse) Reset()         { *m = QueryCustomParamsResponse{} }
func (m *QueryCustomParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryCustomParamsResponse) ProtoMessage()    {}
func (*QueryCustomParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_47cb083fb0607b21, []int{3}
}
func (m *QueryCustomParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCustomParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCustomParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCustomParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCustomParamsResponse.Merge(m, src)
}
func (m *QueryCustomParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryCustomParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCustomParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCustomParamsResponse proto.InternalMessageInfo

func (m *QueryCustomParamsResponse) GetParams() CustomParams {
	if m != nil {
		return m.Params
	}
	return CustomParams{}
}

func init() {
	proto.RegisterType((*QuerySwitchParamsRequest)(nil), "fx.gov.v1.QuerySwitchParamsRequest")
	proto.RegisterType((*QuerySwitchParamsResponse)(nil), "fx.gov.v1.QuerySwitchParamsResponse")
	proto.RegisterType((*QueryCustomParamsRequest)(nil), "fx.gov.v1.QueryCustomParamsRequest")
	proto.RegisterType((*QueryCustomParamsResponse)(nil), "fx.gov.v1.QueryCustomParamsResponse")
}

func init() { proto.RegisterFile("fx/gov/v1/query.proto", fileDescriptor_47cb083fb0607b21) }

var fileDescriptor_47cb083fb0607b21 = []byte{
	// 373 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x4e, 0xea, 0x40,
	0x14, 0x86, 0x5b, 0x72, 0x2f, 0x37, 0xcc, 0x65, 0x35, 0xb9, 0x57, 0xb0, 0x31, 0x95, 0x54, 0x16,
	0x6a, 0x62, 0x27, 0x40, 0x4c, 0x5c, 0xc3, 0x0b, 0x68, 0x8d, 0x1b, 0x37, 0xa4, 0xd4, 0x61, 0x68,
	0x42, 0x7b, 0x4a, 0x67, 0x5a, 0x21, 0xc4, 0x8d, 0x89, 0x3b, 0x17, 0x26, 0xbe, 0x14, 0x4b, 0x12,
	0x37, 0xae, 0x8c, 0x01, 0x1f, 0xc4, 0x30, 0x6d, 0xa0, 0x04, 0x34, 0xee, 0x4e, 0xda, 0xff, 0xfc,
	0xdf, 0x97, 0xd3, 0xa2, 0xff, 0xdd, 0x21, 0x61, 0x10, 0x93, 0xb8, 0x46, 0x06, 0x11, 0x0d, 0x47,
	0x66, 0x10, 0x82, 0x00, 0x5c, 0xe8, 0x0e, 0x4d, 0x06, 0xb1, 0x19, 0xd7, 0xb4, 0x9d, 0x55, 0x22,
	0xb0, 0x43, 0xdb, 0xe3, 0x49, 0x44, 0xfb, 0xc7, 0x80, 0x81, 0x1c, 0xc9, 0x62, 0x4a, 0x9f, 0xee,
	0x31, 0x00, 0xd6, 0xa7, 0xc4, 0x0e, 0x5c, 0x62, 0xfb, 0x3e, 0x08, 0x5b, 0xb8, 0xe0, 0xa7, 0x3b,
	0x86, 0x86, 0xca, 0x17, 0x0b, 0xca, 0xe5, 0xad, 0x2b, 0x9c, 0xde, 0xb9, 0xac, 0xb3, 0xe8, 0x20,
	0xa2, 0x5c, 0x18, 0x16, 0xda, 0xdd, 0xf2, 0x8e, 0x07, 0xe0, 0x73, 0x8a, 0x4f, 0x51, 0x3e, 0x81,
	0x97, 0xd5, 0x8a, 0x7a, 0xf8, 0xb7, 0x5e, 0x32, 0x97, 0x82, 0x66, 0x76, 0xa1, 0xf9, 0x6b, 0xf2,
	0xb6, 0xaf, 0x58, 0x69, 0xd8, 0x68, 0xa4, 0xbc, 0x56, 0xc4, 0x05, 0x78, 0x6b, 0x3c, 0x5c, 0x42,
	0x7f, 0x3c, 0xce, 0xda, 0x51, 0xd8, 0x97, 0x9d, 0x05, 0x2b, 0xef, 0x71, 0x76, 0x15, 0xf6, 0x97,
	0x22, 0xeb, 0x4b, 0x3f, 0x10, 0xc9, 0x2e, 0xac, 0x8b, 0xd4, 0x1f, 0x73, 0xe8, 0xb7, 0x2c, 0xc5,
	0x63, 0x54, 0xcc, 0x0a, 0xe3, 0x83, 0x4c, 0xc1, 0x57, 0xb7, 0xd1, 0xaa, 0xdf, 0x87, 0x12, 0x37,
	0xa3, 0x72, 0xff, 0xf2, 0xf1, 0x9c, 0xd3, 0x70, 0x99, 0xac, 0x3e, 0x19, 0x97, 0xc1, 0x76, 0xa2,
	0x81, 0x1f, 0x54, 0x54, 0xcc, 0x5a, 0x6e, 0xd2, 0xb7, 0x5c, 0x6a, 0x93, 0xbe, 0xed, 0x32, 0xc6,
	0xb1, 0xa4, 0x57, 0xb1, 0x91, 0xa1, 0x3b, 0x32, 0x98, 0xd2, 0xc9, 0x38, 0xbd, 0xf7, 0x5d, 0xb3,
	0x35, 0x99, 0xe9, 0xea, 0x74, 0xa6, 0xab, 0xef, 0x33, 0x5d, 0x7d, 0x9a, 0xeb, 0xca, 0x74, 0xae,
	0x2b, 0xaf, 0x73, 0x5d, 0xb9, 0x3e, 0x62, 0xae, 0xe8, 0x45, 0x1d, 0xd3, 0x01, 0x8f, 0x04, 0x91,
	0x7f, 0xe3, 0xda, 0x2e, 0xe9, 0x0e, 0x4f, 0x1c, 0x08, 0x29, 0x89, 0xcf, 0x48, 0xd2, 0x2c, 0x46,
	0x01, 0xe5, 0x9d, 0xbc, 0xfc, 0xa7, 0x1a, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x02, 0xef, 0x1b,
	0x86, 0xc3, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	SwitchParams(ctx context.Context, in *QuerySwitchParamsRequest, opts ...grpc.CallOption) (*QuerySwitchParamsResponse, error)
	CustomParams(ctx context.Context, in *QueryCustomParamsRequest, opts ...grpc.CallOption) (*QueryCustomParamsResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) SwitchParams(ctx context.Context, in *QuerySwitchParamsRequest, opts ...grpc.CallOption) (*QuerySwitchParamsResponse, error) {
	out := new(QuerySwitchParamsResponse)
	err := c.cc.Invoke(ctx, "/fx.gov.v1.Query/SwitchParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) CustomParams(ctx context.Context, in *QueryCustomParamsRequest, opts ...grpc.CallOption) (*QueryCustomParamsResponse, error) {
	out := new(QueryCustomParamsResponse)
	err := c.cc.Invoke(ctx, "/fx.gov.v1.Query/CustomParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	SwitchParams(context.Context, *QuerySwitchParamsRequest) (*QuerySwitchParamsResponse, error)
	CustomParams(context.Context, *QueryCustomParamsRequest) (*QueryCustomParamsResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) SwitchParams(ctx context.Context, req *QuerySwitchParamsRequest) (*QuerySwitchParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SwitchParams not implemented")
}
func (*UnimplementedQueryServer) CustomParams(ctx context.Context, req *QueryCustomParamsRequest) (*QueryCustomParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CustomParams not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_SwitchParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySwitchParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SwitchParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.gov.v1.Query/SwitchParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SwitchParams(ctx, req.(*QuerySwitchParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CustomParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCustomParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CustomParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.gov.v1.Query/CustomParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CustomParams(ctx, req.(*QueryCustomParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fx.gov.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SwitchParams",
			Handler:    _Query_SwitchParams_Handler,
		},
		{
			MethodName: "CustomParams",
			Handler:    _Query_CustomParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fx/gov/v1/query.proto",
}

func (m *QuerySwitchParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuerySwitchParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuerySwitchParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QuerySwitchParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuerySwitchParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuerySwitchParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryCustomParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCustomParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCustomParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.MsgUrl) > 0 {
		i -= len(m.MsgUrl)
		copy(dAtA[i:], m.MsgUrl)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.MsgUrl)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryCustomParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCustomParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCustomParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QuerySwitchParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QuerySwitchParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryCustomParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.MsgUrl)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryCustomParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QuerySwitchParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuerySwitchParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuerySwitchParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QuerySwitchParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuerySwitchParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuerySwitchParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryCustomParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryCustomParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCustomParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MsgUrl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MsgUrl = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryCustomParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryCustomParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCustomParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
