// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MetricsCollectorGrpcClient is the client API for MetricsCollectorGrpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsCollectorGrpcClient interface {
	Send(ctx context.Context, opts ...grpc.CallOption) (MetricsCollectorGrpc_SendClient, error)
}

type metricsCollectorGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsCollectorGrpcClient(cc grpc.ClientConnInterface) MetricsCollectorGrpcClient {
	return &metricsCollectorGrpcClient{cc}
}

func (c *metricsCollectorGrpcClient) Send(ctx context.Context, opts ...grpc.CallOption) (MetricsCollectorGrpc_SendClient, error) {
	stream, err := c.cc.NewStream(ctx, &MetricsCollectorGrpc_ServiceDesc.Streams[0], "/MetricsCollectorGrpc/Send", opts...)
	if err != nil {
		return nil, err
	}
	x := &metricsCollectorGrpcSendClient{stream}
	return x, nil
}

type MetricsCollectorGrpc_SendClient interface {
	Send(*Metric) error
	Recv() (*empty.Empty, error)
	grpc.ClientStream
}

type metricsCollectorGrpcSendClient struct {
	grpc.ClientStream
}

func (x *metricsCollectorGrpcSendClient) Send(m *Metric) error {
	return x.ClientStream.SendMsg(m)
}

func (x *metricsCollectorGrpcSendClient) Recv() (*empty.Empty, error) {
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MetricsCollectorGrpcServer is the server API for MetricsCollectorGrpc service.
// All implementations must embed UnimplementedMetricsCollectorGrpcServer
// for forward compatibility
type MetricsCollectorGrpcServer interface {
	Send(MetricsCollectorGrpc_SendServer) error
	mustEmbedUnimplementedMetricsCollectorGrpcServer()
}

// UnimplementedMetricsCollectorGrpcServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsCollectorGrpcServer struct {
}

func (UnimplementedMetricsCollectorGrpcServer) Send(MetricsCollectorGrpc_SendServer) error {
	return status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedMetricsCollectorGrpcServer) mustEmbedUnimplementedMetricsCollectorGrpcServer() {}

// UnsafeMetricsCollectorGrpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsCollectorGrpcServer will
// result in compilation errors.
type UnsafeMetricsCollectorGrpcServer interface {
	mustEmbedUnimplementedMetricsCollectorGrpcServer()
}

func RegisterMetricsCollectorGrpcServer(s grpc.ServiceRegistrar, srv MetricsCollectorGrpcServer) {
	s.RegisterService(&MetricsCollectorGrpc_ServiceDesc, srv)
}

func _MetricsCollectorGrpc_Send_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MetricsCollectorGrpcServer).Send(&metricsCollectorGrpcSendServer{stream})
}

type MetricsCollectorGrpc_SendServer interface {
	Send(*empty.Empty) error
	Recv() (*Metric, error)
	grpc.ServerStream
}

type metricsCollectorGrpcSendServer struct {
	grpc.ServerStream
}

func (x *metricsCollectorGrpcSendServer) Send(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *metricsCollectorGrpcSendServer) Recv() (*Metric, error) {
	m := new(Metric)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MetricsCollectorGrpc_ServiceDesc is the grpc.ServiceDesc for MetricsCollectorGrpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricsCollectorGrpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MetricsCollectorGrpc",
	HandlerType: (*MetricsCollectorGrpcServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Send",
			Handler:       _MetricsCollectorGrpc_Send_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "metrics.proto",
}
