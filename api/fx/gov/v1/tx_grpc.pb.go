// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: fx/gov/v1/tx.proto

package govv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Msg_UpdateStore_FullMethodName        = "/fx.gov.v1.Msg/UpdateStore"
	Msg_UpdateSwitchParams_FullMethodName = "/fx.gov.v1.Msg/UpdateSwitchParams"
	Msg_UpdateCustomParams_FullMethodName = "/fx.gov.v1.Msg/UpdateCustomParams"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	UpdateStore(ctx context.Context, in *MsgUpdateStore, opts ...grpc.CallOption) (*MsgUpdateStoreResponse, error)
	UpdateSwitchParams(ctx context.Context, in *MsgUpdateSwitchParams, opts ...grpc.CallOption) (*MsgUpdateSwitchParamsResponse, error)
	UpdateCustomParams(ctx context.Context, in *MsgUpdateCustomParams, opts ...grpc.CallOption) (*MsgUpdateCustomParamsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) UpdateStore(ctx context.Context, in *MsgUpdateStore, opts ...grpc.CallOption) (*MsgUpdateStoreResponse, error) {
	out := new(MsgUpdateStoreResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateStore_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateSwitchParams(ctx context.Context, in *MsgUpdateSwitchParams, opts ...grpc.CallOption) (*MsgUpdateSwitchParamsResponse, error) {
	out := new(MsgUpdateSwitchParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateSwitchParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateCustomParams(ctx context.Context, in *MsgUpdateCustomParams, opts ...grpc.CallOption) (*MsgUpdateCustomParamsResponse, error) {
	out := new(MsgUpdateCustomParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateCustomParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	UpdateStore(context.Context, *MsgUpdateStore) (*MsgUpdateStoreResponse, error)
	UpdateSwitchParams(context.Context, *MsgUpdateSwitchParams) (*MsgUpdateSwitchParamsResponse, error)
	UpdateCustomParams(context.Context, *MsgUpdateCustomParams) (*MsgUpdateCustomParamsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) UpdateStore(context.Context, *MsgUpdateStore) (*MsgUpdateStoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStore not implemented")
}
func (UnimplementedMsgServer) UpdateSwitchParams(context.Context, *MsgUpdateSwitchParams) (*MsgUpdateSwitchParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSwitchParams not implemented")
}
func (UnimplementedMsgServer) UpdateCustomParams(context.Context, *MsgUpdateCustomParams) (*MsgUpdateCustomParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCustomParams not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_UpdateStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateStore)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateStore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateStore(ctx, req.(*MsgUpdateStore))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateSwitchParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateSwitchParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateSwitchParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateSwitchParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateSwitchParams(ctx, req.(*MsgUpdateSwitchParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateCustomParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateCustomParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateCustomParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateCustomParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateCustomParams(ctx, req.(*MsgUpdateCustomParams))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fx.gov.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateStore",
			Handler:    _Msg_UpdateStore_Handler,
		},
		{
			MethodName: "UpdateSwitchParams",
			Handler:    _Msg_UpdateSwitchParams_Handler,
		},
		{
			MethodName: "UpdateCustomParams",
			Handler:    _Msg_UpdateCustomParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fx/gov/v1/tx.proto",
}
