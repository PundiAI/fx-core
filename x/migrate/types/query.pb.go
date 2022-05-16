// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: migrate/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
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

type QueryMigrateRecordRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryMigrateRecordRequest) Reset()         { *m = QueryMigrateRecordRequest{} }
func (m *QueryMigrateRecordRequest) String() string { return proto.CompactTextString(m) }
func (*QueryMigrateRecordRequest) ProtoMessage()    {}
func (*QueryMigrateRecordRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c046ee78a0cdc93, []int{0}
}
func (m *QueryMigrateRecordRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryMigrateRecordRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryMigrateRecordRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryMigrateRecordRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMigrateRecordRequest.Merge(m, src)
}
func (m *QueryMigrateRecordRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryMigrateRecordRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMigrateRecordRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMigrateRecordRequest proto.InternalMessageInfo

func (m *QueryMigrateRecordRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type QueryMigrateRecordResponse struct {
	// has migrate true-> migrated, false-> not migrated.
	Found bool `protobuf:"varint,1,opt,name=found,proto3" json:"found,omitempty"`
	// migrateRecord defines the the migrate record.
	MigrateRecord MigrateRecord `protobuf:"bytes,2,opt,name=migrateRecord,proto3" json:"migrateRecord"`
}

func (m *QueryMigrateRecordResponse) Reset()         { *m = QueryMigrateRecordResponse{} }
func (m *QueryMigrateRecordResponse) String() string { return proto.CompactTextString(m) }
func (*QueryMigrateRecordResponse) ProtoMessage()    {}
func (*QueryMigrateRecordResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c046ee78a0cdc93, []int{1}
}
func (m *QueryMigrateRecordResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryMigrateRecordResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryMigrateRecordResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryMigrateRecordResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMigrateRecordResponse.Merge(m, src)
}
func (m *QueryMigrateRecordResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryMigrateRecordResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMigrateRecordResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMigrateRecordResponse proto.InternalMessageInfo

func (m *QueryMigrateRecordResponse) GetFound() bool {
	if m != nil {
		return m.Found
	}
	return false
}

func (m *QueryMigrateRecordResponse) GetMigrateRecord() MigrateRecord {
	if m != nil {
		return m.MigrateRecord
	}
	return MigrateRecord{}
}

type QueryMigrateCheckAccountRequest struct {
	// migrate from address
	From string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	// migrate to address
	To string `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
}

func (m *QueryMigrateCheckAccountRequest) Reset()         { *m = QueryMigrateCheckAccountRequest{} }
func (m *QueryMigrateCheckAccountRequest) String() string { return proto.CompactTextString(m) }
func (*QueryMigrateCheckAccountRequest) ProtoMessage()    {}
func (*QueryMigrateCheckAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c046ee78a0cdc93, []int{2}
}
func (m *QueryMigrateCheckAccountRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryMigrateCheckAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryMigrateCheckAccountRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryMigrateCheckAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMigrateCheckAccountRequest.Merge(m, src)
}
func (m *QueryMigrateCheckAccountRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryMigrateCheckAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMigrateCheckAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMigrateCheckAccountRequest proto.InternalMessageInfo

func (m *QueryMigrateCheckAccountRequest) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *QueryMigrateCheckAccountRequest) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

type QueryMigrateCheckAccountResponse struct {
}

func (m *QueryMigrateCheckAccountResponse) Reset()         { *m = QueryMigrateCheckAccountResponse{} }
func (m *QueryMigrateCheckAccountResponse) String() string { return proto.CompactTextString(m) }
func (*QueryMigrateCheckAccountResponse) ProtoMessage()    {}
func (*QueryMigrateCheckAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c046ee78a0cdc93, []int{3}
}
func (m *QueryMigrateCheckAccountResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryMigrateCheckAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryMigrateCheckAccountResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryMigrateCheckAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMigrateCheckAccountResponse.Merge(m, src)
}
func (m *QueryMigrateCheckAccountResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryMigrateCheckAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMigrateCheckAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMigrateCheckAccountResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryMigrateRecordRequest)(nil), "fx.migrate.v1.QueryMigrateRecordRequest")
	proto.RegisterType((*QueryMigrateRecordResponse)(nil), "fx.migrate.v1.QueryMigrateRecordResponse")
	proto.RegisterType((*QueryMigrateCheckAccountRequest)(nil), "fx.migrate.v1.QueryMigrateCheckAccountRequest")
	proto.RegisterType((*QueryMigrateCheckAccountResponse)(nil), "fx.migrate.v1.QueryMigrateCheckAccountResponse")
}

func init() { proto.RegisterFile("migrate/v1/query.proto", fileDescriptor_7c046ee78a0cdc93) }

var fileDescriptor_7c046ee78a0cdc93 = []byte{
	// 410 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xb1, 0xae, 0xd3, 0x30,
	0x14, 0x86, 0xe3, 0xe8, 0x5e, 0xa0, 0x46, 0x65, 0x30, 0x15, 0x4a, 0x43, 0x95, 0x96, 0x88, 0xa1,
	0x48, 0x10, 0xab, 0x45, 0x3c, 0x00, 0x45, 0x20, 0x16, 0x06, 0x32, 0xb2, 0xa5, 0x8e, 0x93, 0x46,
	0x90, 0x9c, 0xd4, 0x71, 0xaa, 0x54, 0xc0, 0xc2, 0xc8, 0x84, 0xc4, 0xc4, 0x5b, 0xf0, 0x18, 0x1d,
	0x2b, 0xb1, 0x30, 0x21, 0xd4, 0xf2, 0x20, 0x28, 0x4e, 0x8a, 0x12, 0xd4, 0x5e, 0x75, 0x3b, 0xb6,
	0xff, 0xf3, 0xff, 0x9f, 0x8f, 0x8d, 0xef, 0xc4, 0x51, 0x28, 0x3c, 0xc9, 0xe9, 0x6a, 0x42, 0x97,
	0x39, 0x17, 0x6b, 0x27, 0x15, 0x20, 0x81, 0x74, 0x83, 0xc2, 0xa9, 0x8f, 0x9c, 0xd5, 0xc4, 0x1c,
	0x84, 0x00, 0xe1, 0x3b, 0x4e, 0xbd, 0x34, 0xa2, 0x5e, 0x92, 0x80, 0xf4, 0x64, 0x04, 0x49, 0x56,
	0x89, 0xcd, 0x5e, 0x08, 0x21, 0xa8, 0x92, 0x96, 0x55, 0xbd, 0x6b, 0x34, 0xac, 0x0f, 0x56, 0xea,
	0xc4, 0x7e, 0x82, 0xfb, 0xaf, 0xcb, 0xac, 0x57, 0xd5, 0xae, 0xcb, 0x19, 0x08, 0xdf, 0xe5, 0xcb,
	0x9c, 0x67, 0x92, 0x18, 0xf8, 0xba, 0xe7, 0xfb, 0x82, 0x67, 0x99, 0x81, 0x46, 0x68, 0xdc, 0x71,
	0x0f, 0x4b, 0xfb, 0x03, 0x36, 0x8f, 0xb5, 0x65, 0x29, 0x24, 0x19, 0x27, 0x3d, 0x7c, 0x19, 0x40,
	0x9e, 0xf8, 0xaa, 0xeb, 0x86, 0x5b, 0x2d, 0xc8, 0x4b, 0xdc, 0x8d, 0x9b, 0x72, 0x43, 0x1f, 0xa1,
	0xf1, 0xcd, 0xe9, 0xc0, 0x69, 0xdd, 0xcf, 0x69, 0x59, 0xce, 0x2e, 0x36, 0xbf, 0x86, 0x9a, 0xdb,
	0x6e, 0xb4, 0x9f, 0xe3, 0x61, 0x33, 0xfd, 0xd9, 0x82, 0xb3, 0xb7, 0x4f, 0x19, 0x83, 0x3c, 0x91,
	0x07, 0x74, 0x82, 0x2f, 0x02, 0x01, 0x71, 0xcd, 0xad, 0x6a, 0x72, 0x0b, 0xeb, 0x12, 0x54, 0x6a,
	0xc7, 0xd5, 0x25, 0xd8, 0x36, 0x1e, 0x9d, 0xb6, 0xa9, 0xae, 0x32, 0xfd, 0xae, 0xe3, 0x4b, 0x25,
	0x22, 0x9f, 0x11, 0xee, 0xb6, 0xd8, 0xc8, 0xf8, 0x3f, 0xf2, 0x93, 0x83, 0x34, 0x1f, 0x9c, 0xa1,
	0xac, 0x02, 0xed, 0xfb, 0x9f, 0x7e, 0xfc, 0xf9, 0xaa, 0x5b, 0x64, 0x40, 0x1b, 0x6f, 0x26, 0x94,
	0x86, 0xbe, 0xaf, 0xc7, 0xff, 0x91, 0x7c, 0x43, 0xf8, 0xf6, 0x11, 0x6c, 0xe2, 0x5c, 0x11, 0x74,
	0x64, 0x4c, 0x26, 0x3d, 0x5b, 0x5f, 0xe3, 0xdd, 0x53, 0x78, 0x77, 0x49, 0xbf, 0x89, 0xc7, 0x4a,
	0x25, 0xf5, 0x2a, 0xe9, 0xec, 0xc5, 0x66, 0x67, 0xa1, 0xed, 0xce, 0x42, 0xbf, 0x77, 0x16, 0xfa,
	0xb2, 0xb7, 0xb4, 0xed, 0xde, 0xd2, 0x7e, 0xee, 0x2d, 0xed, 0xcd, 0xc3, 0x30, 0x92, 0x8b, 0x7c,
	0xee, 0x30, 0x88, 0x69, 0x90, 0x27, 0xac, 0xfc, 0xb8, 0x05, 0x0d, 0x8a, 0x47, 0x0c, 0x04, 0xa7,
	0xc5, 0x3f, 0x4b, 0xb9, 0x4e, 0x79, 0x36, 0xbf, 0xa6, 0x7e, 0xe8, 0xe3, 0xbf, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x61, 0x6f, 0xf2, 0xbc, 0x18, 0x03, 0x00, 0x00,
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
	// DenomTrace queries a denomination trace information.
	MigrateRecord(ctx context.Context, in *QueryMigrateRecordRequest, opts ...grpc.CallOption) (*QueryMigrateRecordResponse, error)
	MigrateCheckAccount(ctx context.Context, in *QueryMigrateCheckAccountRequest, opts ...grpc.CallOption) (*QueryMigrateCheckAccountResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) MigrateRecord(ctx context.Context, in *QueryMigrateRecordRequest, opts ...grpc.CallOption) (*QueryMigrateRecordResponse, error) {
	out := new(QueryMigrateRecordResponse)
	err := c.cc.Invoke(ctx, "/fx.migrate.v1.Query/MigrateRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) MigrateCheckAccount(ctx context.Context, in *QueryMigrateCheckAccountRequest, opts ...grpc.CallOption) (*QueryMigrateCheckAccountResponse, error) {
	out := new(QueryMigrateCheckAccountResponse)
	err := c.cc.Invoke(ctx, "/fx.migrate.v1.Query/MigrateCheckAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// DenomTrace queries a denomination trace information.
	MigrateRecord(context.Context, *QueryMigrateRecordRequest) (*QueryMigrateRecordResponse, error)
	MigrateCheckAccount(context.Context, *QueryMigrateCheckAccountRequest) (*QueryMigrateCheckAccountResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) MigrateRecord(ctx context.Context, req *QueryMigrateRecordRequest) (*QueryMigrateRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MigrateRecord not implemented")
}
func (*UnimplementedQueryServer) MigrateCheckAccount(ctx context.Context, req *QueryMigrateCheckAccountRequest) (*QueryMigrateCheckAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MigrateCheckAccount not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_MigrateRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryMigrateRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).MigrateRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.migrate.v1.Query/MigrateRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).MigrateRecord(ctx, req.(*QueryMigrateRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_MigrateCheckAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryMigrateCheckAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).MigrateCheckAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.migrate.v1.Query/MigrateCheckAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).MigrateCheckAccount(ctx, req.(*QueryMigrateCheckAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fx.migrate.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MigrateRecord",
			Handler:    _Query_MigrateRecord_Handler,
		},
		{
			MethodName: "MigrateCheckAccount",
			Handler:    _Query_MigrateCheckAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "migrate/v1/query.proto",
}

func (m *QueryMigrateRecordRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryMigrateRecordRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryMigrateRecordRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryMigrateRecordResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryMigrateRecordResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryMigrateRecordResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.MigrateRecord.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.Found {
		i--
		if m.Found {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryMigrateCheckAccountRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryMigrateCheckAccountRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryMigrateCheckAccountRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.To) > 0 {
		i -= len(m.To)
		copy(dAtA[i:], m.To)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.To)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.From) > 0 {
		i -= len(m.From)
		copy(dAtA[i:], m.From)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.From)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryMigrateCheckAccountResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryMigrateCheckAccountResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryMigrateCheckAccountResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
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
func (m *QueryMigrateRecordRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryMigrateRecordResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Found {
		n += 2
	}
	l = m.MigrateRecord.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryMigrateCheckAccountRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.From)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.To)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryMigrateCheckAccountResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryMigrateRecordRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryMigrateRecordRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryMigrateRecordRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			m.Address = string(dAtA[iNdEx:postIndex])
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
func (m *QueryMigrateRecordResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryMigrateRecordResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryMigrateRecordResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Found", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Found = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MigrateRecord", wireType)
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
			if err := m.MigrateRecord.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryMigrateCheckAccountRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryMigrateCheckAccountRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryMigrateCheckAccountRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
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
			m.From = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field To", wireType)
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
			m.To = string(dAtA[iNdEx:postIndex])
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
func (m *QueryMigrateCheckAccountResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryMigrateCheckAccountResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryMigrateCheckAccountResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
