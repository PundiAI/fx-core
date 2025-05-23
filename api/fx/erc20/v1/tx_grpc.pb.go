// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: fx/erc20/v1/tx.proto

package erc20v1

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
	Msg_ConvertCoin_FullMethodName           = "/fx.erc20.v1.Msg/ConvertCoin"
	Msg_UpdateParams_FullMethodName          = "/fx.erc20.v1.Msg/UpdateParams"
	Msg_ToggleTokenConversion_FullMethodName = "/fx.erc20.v1.Msg/ToggleTokenConversion"
	Msg_RegisterNativeCoin_FullMethodName    = "/fx.erc20.v1.Msg/RegisterNativeCoin"
	Msg_RegisterNativeERC20_FullMethodName   = "/fx.erc20.v1.Msg/RegisterNativeERC20"
	Msg_RegisterBridgeToken_FullMethodName   = "/fx.erc20.v1.Msg/RegisterBridgeToken"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// ConvertCoin mints a ERC20 representation of the SDK Coin denom that is
	// registered on the token mapping.
	ConvertCoin(ctx context.Context, in *MsgConvertCoin, opts ...grpc.CallOption) (*MsgConvertCoinResponse, error)
	// UpdateParams defines a governance operation for updating the x/erc20 module
	// parameters. The authority is hard-coded to the x/gov module account.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	ToggleTokenConversion(ctx context.Context, in *MsgToggleTokenConversion, opts ...grpc.CallOption) (*MsgToggleTokenConversionResponse, error)
	RegisterNativeCoin(ctx context.Context, in *MsgRegisterNativeCoin, opts ...grpc.CallOption) (*MsgRegisterNativeCoinResponse, error)
	RegisterNativeERC20(ctx context.Context, in *MsgRegisterNativeERC20, opts ...grpc.CallOption) (*MsgRegisterNativeERC20Response, error)
	RegisterBridgeToken(ctx context.Context, in *MsgRegisterBridgeToken, opts ...grpc.CallOption) (*MsgRegisterBridgeTokenResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) ConvertCoin(ctx context.Context, in *MsgConvertCoin, opts ...grpc.CallOption) (*MsgConvertCoinResponse, error) {
	out := new(MsgConvertCoinResponse)
	err := c.cc.Invoke(ctx, Msg_ConvertCoin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ToggleTokenConversion(ctx context.Context, in *MsgToggleTokenConversion, opts ...grpc.CallOption) (*MsgToggleTokenConversionResponse, error) {
	out := new(MsgToggleTokenConversionResponse)
	err := c.cc.Invoke(ctx, Msg_ToggleTokenConversion_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RegisterNativeCoin(ctx context.Context, in *MsgRegisterNativeCoin, opts ...grpc.CallOption) (*MsgRegisterNativeCoinResponse, error) {
	out := new(MsgRegisterNativeCoinResponse)
	err := c.cc.Invoke(ctx, Msg_RegisterNativeCoin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RegisterNativeERC20(ctx context.Context, in *MsgRegisterNativeERC20, opts ...grpc.CallOption) (*MsgRegisterNativeERC20Response, error) {
	out := new(MsgRegisterNativeERC20Response)
	err := c.cc.Invoke(ctx, Msg_RegisterNativeERC20_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RegisterBridgeToken(ctx context.Context, in *MsgRegisterBridgeToken, opts ...grpc.CallOption) (*MsgRegisterBridgeTokenResponse, error) {
	out := new(MsgRegisterBridgeTokenResponse)
	err := c.cc.Invoke(ctx, Msg_RegisterBridgeToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// ConvertCoin mints a ERC20 representation of the SDK Coin denom that is
	// registered on the token mapping.
	ConvertCoin(context.Context, *MsgConvertCoin) (*MsgConvertCoinResponse, error)
	// UpdateParams defines a governance operation for updating the x/erc20 module
	// parameters. The authority is hard-coded to the x/gov module account.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	ToggleTokenConversion(context.Context, *MsgToggleTokenConversion) (*MsgToggleTokenConversionResponse, error)
	RegisterNativeCoin(context.Context, *MsgRegisterNativeCoin) (*MsgRegisterNativeCoinResponse, error)
	RegisterNativeERC20(context.Context, *MsgRegisterNativeERC20) (*MsgRegisterNativeERC20Response, error)
	RegisterBridgeToken(context.Context, *MsgRegisterBridgeToken) (*MsgRegisterBridgeTokenResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) ConvertCoin(context.Context, *MsgConvertCoin) (*MsgConvertCoinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConvertCoin not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) ToggleTokenConversion(context.Context, *MsgToggleTokenConversion) (*MsgToggleTokenConversionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToggleTokenConversion not implemented")
}
func (UnimplementedMsgServer) RegisterNativeCoin(context.Context, *MsgRegisterNativeCoin) (*MsgRegisterNativeCoinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterNativeCoin not implemented")
}
func (UnimplementedMsgServer) RegisterNativeERC20(context.Context, *MsgRegisterNativeERC20) (*MsgRegisterNativeERC20Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterNativeERC20 not implemented")
}
func (UnimplementedMsgServer) RegisterBridgeToken(context.Context, *MsgRegisterBridgeToken) (*MsgRegisterBridgeTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterBridgeToken not implemented")
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

func _Msg_ConvertCoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgConvertCoin)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ConvertCoin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_ConvertCoin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ConvertCoin(ctx, req.(*MsgConvertCoin))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ToggleTokenConversion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgToggleTokenConversion)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ToggleTokenConversion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_ToggleTokenConversion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ToggleTokenConversion(ctx, req.(*MsgToggleTokenConversion))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RegisterNativeCoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterNativeCoin)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterNativeCoin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RegisterNativeCoin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterNativeCoin(ctx, req.(*MsgRegisterNativeCoin))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RegisterNativeERC20_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterNativeERC20)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterNativeERC20(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RegisterNativeERC20_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterNativeERC20(ctx, req.(*MsgRegisterNativeERC20))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RegisterBridgeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterBridgeToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterBridgeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RegisterBridgeToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterBridgeToken(ctx, req.(*MsgRegisterBridgeToken))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fx.erc20.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConvertCoin",
			Handler:    _Msg_ConvertCoin_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
		{
			MethodName: "ToggleTokenConversion",
			Handler:    _Msg_ToggleTokenConversion_Handler,
		},
		{
			MethodName: "RegisterNativeCoin",
			Handler:    _Msg_RegisterNativeCoin_Handler,
		},
		{
			MethodName: "RegisterNativeERC20",
			Handler:    _Msg_RegisterNativeERC20_Handler,
		},
		{
			MethodName: "RegisterBridgeToken",
			Handler:    _Msg_RegisterBridgeToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fx/erc20/v1/tx.proto",
}
