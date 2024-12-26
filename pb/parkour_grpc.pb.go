// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.19.4
// source: parkour.proto

package pb

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
	Parkour_Ping_FullMethodName = "/pb.Parkour/Ping"
)

// ParkourClient is the client API for Parkour service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ParkourClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
}

type parkourClient struct {
	cc grpc.ClientConnInterface
}

func NewParkourClient(cc grpc.ClientConnInterface) ParkourClient {
	return &parkourClient{cc}
}

func (c *parkourClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, Parkour_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ParkourServer is the server API for Parkour service.
// All implementations must embed UnimplementedParkourServer
// for forward compatibility
type ParkourServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	mustEmbedUnimplementedParkourServer()
}

// UnimplementedParkourServer must be embedded to have forward compatible implementations.
type UnimplementedParkourServer struct {
}

func (UnimplementedParkourServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedParkourServer) mustEmbedUnimplementedParkourServer() {}

// UnsafeParkourServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ParkourServer will
// result in compilation errors.
type UnsafeParkourServer interface {
	mustEmbedUnimplementedParkourServer()
}

func RegisterParkourServer(s grpc.ServiceRegistrar, srv ParkourServer) {
	s.RegisterService(&Parkour_ServiceDesc, srv)
}

func _Parkour_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ParkourServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Parkour_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ParkourServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Parkour_ServiceDesc is the grpc.ServiceDesc for Parkour service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Parkour_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Parkour",
	HandlerType: (*ParkourServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Parkour_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "parkour.proto",
}
