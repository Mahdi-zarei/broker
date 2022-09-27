// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.3
// source: broker.proto

package srv

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

// BrokerClient is the client API for Broker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrokerClient interface {
	// Publish returns an id if the delivery is successful
	// If broker is closed, should return Unavailable
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	// Subscribe returns an stream of messages
	// If broker is closed, should return Unavailable
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Broker_SubscribeClient, error)
	// Fetch returns the proper message body, if its present
	// If broker is closed, should return Unavailable
	// If the provided id is expired or not present,
	// should return InvalidArgument
	Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*MessageResponse, error)
}

type brokerClient struct {
	cc grpc.ClientConnInterface
}

func NewBrokerClient(cc grpc.ClientConnInterface) BrokerClient {
	return &brokerClient{cc}
}

func (c *brokerClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, "/broker.Broker/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Broker_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Broker_ServiceDesc.Streams[0], "/broker.Broker/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &brokerSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Broker_SubscribeClient interface {
	Recv() (*MessageResponse, error)
	grpc.ClientStream
}

type brokerSubscribeClient struct {
	grpc.ClientStream
}

func (x *brokerSubscribeClient) Recv() (*MessageResponse, error) {
	m := new(MessageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *brokerClient) Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, "/broker.Broker/Fetch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BrokerServer is the server API for Broker service.
// All implementations must embed UnimplementedBrokerServer
// for forward compatibility
type BrokerServer interface {
	// Publish returns an id if the delivery is successful
	// If broker is closed, should return Unavailable
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	// Subscribe returns an stream of messages
	// If broker is closed, should return Unavailable
	Subscribe(*SubscribeRequest, Broker_SubscribeServer) error
	// Fetch returns the proper message body, if its present
	// If broker is closed, should return Unavailable
	// If the provided id is expired or not present,
	// should return InvalidArgument
	Fetch(context.Context, *FetchRequest) (*MessageResponse, error)
	mustEmbedUnimplementedBrokerServer()
}

// UnimplementedBrokerServer must be embedded to have forward compatible implementations.
type UnimplementedBrokerServer struct {
}

func (UnimplementedBrokerServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedBrokerServer) Subscribe(*SubscribeRequest, Broker_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedBrokerServer) Fetch(context.Context, *FetchRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fetch not implemented")
}
func (UnimplementedBrokerServer) mustEmbedUnimplementedBrokerServer() {}

// UnsafeBrokerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrokerServer will
// result in compilation errors.
type UnsafeBrokerServer interface {
	mustEmbedUnimplementedBrokerServer()
}

func RegisterBrokerServer(s grpc.ServiceRegistrar, srv BrokerServer) {
	s.RegisterService(&Broker_ServiceDesc, srv)
}

func _Broker_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/broker.Broker/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BrokerServer).Subscribe(m, &brokerSubscribeServer{stream})
}

type Broker_SubscribeServer interface {
	Send(*MessageResponse) error
	grpc.ServerStream
}

type brokerSubscribeServer struct {
	grpc.ServerStream
}

func (x *brokerSubscribeServer) Send(m *MessageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Broker_Fetch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).Fetch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/broker.Broker/Fetch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).Fetch(ctx, req.(*FetchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Broker_ServiceDesc is the grpc.ServiceDesc for Broker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "broker.Broker",
	HandlerType: (*BrokerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _Broker_Publish_Handler,
		},
		{
			MethodName: "Fetch",
			Handler:    _Broker_Fetch_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Broker_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "broker.proto",
}
