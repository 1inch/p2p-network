// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.6
// source: resolver.proto

package resolver

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Execute_Execute_FullMethodName = "/resolver.Execute/Execute"
)

// ExecuteClient is the client API for Execute service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExecuteClient interface {
	Execute(ctx context.Context, in *ResolverRequest, opts ...grpc.CallOption) (*ResolverResponse, error)
}

type executeClient struct {
	cc grpc.ClientConnInterface
}

func NewExecuteClient(cc grpc.ClientConnInterface) ExecuteClient {
	return &executeClient{cc}
}

func (c *executeClient) Execute(ctx context.Context, in *ResolverRequest, opts ...grpc.CallOption) (*ResolverResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResolverResponse)
	err := c.cc.Invoke(ctx, Execute_Execute_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExecuteServer is the server API for Execute service.
// All implementations must embed UnimplementedExecuteServer
// for forward compatibility.
type ExecuteServer interface {
	Execute(context.Context, *ResolverRequest) (*ResolverResponse, error)
	mustEmbedUnimplementedExecuteServer()
}

// UnimplementedExecuteServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedExecuteServer struct{}

func (UnimplementedExecuteServer) Execute(context.Context, *ResolverRequest) (*ResolverResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedExecuteServer) mustEmbedUnimplementedExecuteServer() {}
func (UnimplementedExecuteServer) testEmbeddedByValue()                 {}

// UnsafeExecuteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExecuteServer will
// result in compilation errors.
type UnsafeExecuteServer interface {
	mustEmbedUnimplementedExecuteServer()
}

func RegisterExecuteServer(s grpc.ServiceRegistrar, srv ExecuteServer) {
	// If the following call pancis, it indicates UnimplementedExecuteServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Execute_ServiceDesc, srv)
}

func _Execute_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResolverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecuteServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Execute_Execute_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecuteServer).Execute(ctx, req.(*ResolverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Execute_ServiceDesc is the grpc.ServiceDesc for Execute service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Execute_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "resolver.Execute",
	HandlerType: (*ExecuteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Execute",
			Handler:    _Execute_Execute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "resolver.proto",
}
