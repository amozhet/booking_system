// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.0
// source: booking_system/roomManage/proto/room.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	RoomService_GetRooms_FullMethodName    = "/proto.RoomService/GetRooms"
	RoomService_GetRoomByID_FullMethodName = "/proto.RoomService/GetRoomByID"
	RoomService_CreateRoom_FullMethodName  = "/proto.RoomService/CreateRoom"
	RoomService_UpdateRoom_FullMethodName  = "/proto.RoomService/UpdateRoom"
	RoomService_DeleteRoom_FullMethodName  = "/proto.RoomService/DeleteRoom"
)

// RoomServiceClient is the client API for RoomService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomServiceClient interface {
	GetRooms(ctx context.Context, in *GetRoomsRequest, opts ...grpc.CallOption) (*GetRoomsResponse, error)
	GetRoomByID(ctx context.Context, in *GetRoomByIDRequest, opts ...grpc.CallOption) (*GetRoomResponse, error)
	CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*RoomResponse, error)
	UpdateRoom(ctx context.Context, in *UpdateRoomRequest, opts ...grpc.CallOption) (*RoomResponse, error)
	DeleteRoom(ctx context.Context, in *DeleteRoomRequest, opts ...grpc.CallOption) (*DeleteRoomResponse, error)
}

type roomServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomServiceClient(cc grpc.ClientConnInterface) RoomServiceClient {
	return &roomServiceClient{cc}
}

func (c *roomServiceClient) GetRooms(ctx context.Context, in *GetRoomsRequest, opts ...grpc.CallOption) (*GetRoomsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomsResponse)
	err := c.cc.Invoke(ctx, RoomService_GetRooms_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) GetRoomByID(ctx context.Context, in *GetRoomByIDRequest, opts ...grpc.CallOption) (*GetRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomResponse)
	err := c.cc.Invoke(ctx, RoomService_GetRoomByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*RoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RoomResponse)
	err := c.cc.Invoke(ctx, RoomService_CreateRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) UpdateRoom(ctx context.Context, in *UpdateRoomRequest, opts ...grpc.CallOption) (*RoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RoomResponse)
	err := c.cc.Invoke(ctx, RoomService_UpdateRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) DeleteRoom(ctx context.Context, in *DeleteRoomRequest, opts ...grpc.CallOption) (*DeleteRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteRoomResponse)
	err := c.cc.Invoke(ctx, RoomService_DeleteRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomServiceServer is the server API for RoomService service.
// All implementations must embed UnimplementedRoomServiceServer
// for forward compatibility
type RoomServiceServer interface {
	GetRooms(context.Context, *GetRoomsRequest) (*GetRoomsResponse, error)
	GetRoomByID(context.Context, *GetRoomByIDRequest) (*GetRoomResponse, error)
	CreateRoom(context.Context, *CreateRoomRequest) (*RoomResponse, error)
	UpdateRoom(context.Context, *UpdateRoomRequest) (*RoomResponse, error)
	DeleteRoom(context.Context, *DeleteRoomRequest) (*DeleteRoomResponse, error)
	mustEmbedUnimplementedRoomServiceServer()
}

// UnimplementedRoomServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRoomServiceServer struct {
}

func (UnimplementedRoomServiceServer) GetRooms(context.Context, *GetRoomsRequest) (*GetRoomsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRooms not implemented")
}
func (UnimplementedRoomServiceServer) GetRoomByID(context.Context, *GetRoomByIDRequest) (*GetRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomByID not implemented")
}
func (UnimplementedRoomServiceServer) CreateRoom(context.Context, *CreateRoomRequest) (*RoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedRoomServiceServer) UpdateRoom(context.Context, *UpdateRoomRequest) (*RoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRoom not implemented")
}
func (UnimplementedRoomServiceServer) DeleteRoom(context.Context, *DeleteRoomRequest) (*DeleteRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRoom not implemented")
}
func (UnimplementedRoomServiceServer) mustEmbedUnimplementedRoomServiceServer() {}

// UnsafeRoomServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomServiceServer will
// result in compilation errors.
type UnsafeRoomServiceServer interface {
	mustEmbedUnimplementedRoomServiceServer()
}

func RegisterRoomServiceServer(s grpc.ServiceRegistrar, srv RoomServiceServer) {
	s.RegisterService(&RoomService_ServiceDesc, srv)
}

func _RoomService_GetRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).GetRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_GetRooms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).GetRooms(ctx, req.(*GetRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_GetRoomByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).GetRoomByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_GetRoomByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).GetRoomByID(ctx, req.(*GetRoomByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_CreateRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).CreateRoom(ctx, req.(*CreateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_UpdateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).UpdateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_UpdateRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).UpdateRoom(ctx, req.(*UpdateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_DeleteRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).DeleteRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_DeleteRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).DeleteRoom(ctx, req.(*DeleteRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoomService_ServiceDesc is the grpc.ServiceDesc for RoomService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoomService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.RoomService",
	HandlerType: (*RoomServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRooms",
			Handler:    _RoomService_GetRooms_Handler,
		},
		{
			MethodName: "GetRoomByID",
			Handler:    _RoomService_GetRoomByID_Handler,
		},
		{
			MethodName: "CreateRoom",
			Handler:    _RoomService_CreateRoom_Handler,
		},
		{
			MethodName: "UpdateRoom",
			Handler:    _RoomService_UpdateRoom_Handler,
		},
		{
			MethodName: "DeleteRoom",
			Handler:    _RoomService_DeleteRoom_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "booking_system/roomManage/proto/room.proto",
}
