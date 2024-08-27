// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.18.1
// source: session-manager.proto

package manager

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SessionManagerClient is the client API for SessionManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SessionManagerClient interface {
	OnSessionCreated(ctx context.Context, in *SessionCreatedRequest, opts ...grpc.CallOption) (*SessionResponse, error)
	OnSessionUpdated(ctx context.Context, in *SessionUpdatedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	OnSessionDeleted(ctx context.Context, in *SessionDeletedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type sessionManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewSessionManagerClient(cc grpc.ClientConnInterface) SessionManagerClient {
	return &sessionManagerClient{cc}
}

func (c *sessionManagerClient) OnSessionCreated(ctx context.Context, in *SessionCreatedRequest, opts ...grpc.CallOption) (*SessionResponse, error) {
	out := new(SessionResponse)
	err := c.cc.Invoke(ctx, "/accelbyte.session.manager.SessionManager/OnSessionCreated", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionManagerClient) OnSessionUpdated(ctx context.Context, in *SessionUpdatedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/accelbyte.session.manager.SessionManager/OnSessionUpdated", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionManagerClient) OnSessionDeleted(ctx context.Context, in *SessionDeletedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/accelbyte.session.manager.SessionManager/OnSessionDeleted", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SessionManagerServer is the server API for SessionManager service.
// All implementations must embed UnimplementedSessionManagerServer
// for forward compatibility
type SessionManagerServer interface {
	OnSessionCreated(context.Context, *SessionCreatedRequest) (*SessionResponse, error)
	OnSessionUpdated(context.Context, *SessionUpdatedRequest) (*emptypb.Empty, error)
	OnSessionDeleted(context.Context, *SessionDeletedRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedSessionManagerServer()
}

// UnimplementedSessionManagerServer must be embedded to have forward compatible implementations.
type UnimplementedSessionManagerServer struct {
}

func (UnimplementedSessionManagerServer) OnSessionCreated(context.Context, *SessionCreatedRequest) (*SessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnSessionCreated not implemented")
}
func (UnimplementedSessionManagerServer) OnSessionUpdated(context.Context, *SessionUpdatedRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnSessionUpdated not implemented")
}
func (UnimplementedSessionManagerServer) OnSessionDeleted(context.Context, *SessionDeletedRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnSessionDeleted not implemented")
}
func (UnimplementedSessionManagerServer) mustEmbedUnimplementedSessionManagerServer() {}

// UnsafeSessionManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SessionManagerServer will
// result in compilation errors.
type UnsafeSessionManagerServer interface {
	mustEmbedUnimplementedSessionManagerServer()
}

func RegisterSessionManagerServer(s grpc.ServiceRegistrar, srv SessionManagerServer) {
	s.RegisterService(&SessionManager_ServiceDesc, srv)
}

func _SessionManager_OnSessionCreated_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionCreatedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).OnSessionCreated(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accelbyte.session.manager.SessionManager/OnSessionCreated",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).OnSessionCreated(ctx, req.(*SessionCreatedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionManager_OnSessionUpdated_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionUpdatedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).OnSessionUpdated(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accelbyte.session.manager.SessionManager/OnSessionUpdated",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).OnSessionUpdated(ctx, req.(*SessionUpdatedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SessionManager_OnSessionDeleted_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionDeletedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionManagerServer).OnSessionDeleted(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accelbyte.session.manager.SessionManager/OnSessionDeleted",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionManagerServer).OnSessionDeleted(ctx, req.(*SessionDeletedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SessionManager_ServiceDesc is the grpc.ServiceDesc for SessionManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SessionManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "accelbyte.session.manager.SessionManager",
	HandlerType: (*SessionManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OnSessionCreated",
			Handler:    _SessionManager_OnSessionCreated_Handler,
		},
		{
			MethodName: "OnSessionUpdated",
			Handler:    _SessionManager_OnSessionUpdated_Handler,
		},
		{
			MethodName: "OnSessionDeleted",
			Handler:    _SessionManager_OnSessionDeleted_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "session-manager.proto",
}
