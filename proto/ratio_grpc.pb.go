// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CurrencyFetcherClient is the client API for CurrencyFetcher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CurrencyFetcherClient interface {
	GetRatios(ctx context.Context, in *GetRatiosRequest, opts ...grpc.CallOption) (*GetRatiosResponse, error)
}

type currencyFetcherClient struct {
	cc grpc.ClientConnInterface
}

func NewCurrencyFetcherClient(cc grpc.ClientConnInterface) CurrencyFetcherClient {
	return &currencyFetcherClient{cc}
}

func (c *currencyFetcherClient) GetRatios(ctx context.Context, in *GetRatiosRequest, opts ...grpc.CallOption) (*GetRatiosResponse, error) {
	out := new(GetRatiosResponse)
	err := c.cc.Invoke(ctx, "/CurrencyFetcher/GetRatios", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CurrencyFetcherServer is the server API for CurrencyFetcher service.
// All implementations must embed UnimplementedCurrencyFetcherServer
// for forward compatibility
type CurrencyFetcherServer interface {
	GetRatios(context.Context, *GetRatiosRequest) (*GetRatiosResponse, error)
	mustEmbedUnimplementedCurrencyFetcherServer()
}

// UnimplementedCurrencyFetcherServer must be embedded to have forward compatible implementations.
type UnimplementedCurrencyFetcherServer struct {
}

func (*UnimplementedCurrencyFetcherServer) GetRatios(context.Context, *GetRatiosRequest) (*GetRatiosResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRatios not implemented")
}
func (*UnimplementedCurrencyFetcherServer) mustEmbedUnimplementedCurrencyFetcherServer() {}

func RegisterCurrencyFetcherServer(s *grpc.Server, srv CurrencyFetcherServer) {
	s.RegisterService(&_CurrencyFetcher_serviceDesc, srv)
}

func _CurrencyFetcher_GetRatios_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRatiosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CurrencyFetcherServer).GetRatios(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CurrencyFetcher/GetRatios",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CurrencyFetcherServer).GetRatios(ctx, req.(*GetRatiosRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CurrencyFetcher_serviceDesc = grpc.ServiceDesc{
	ServiceName: "CurrencyFetcher",
	HandlerType: (*CurrencyFetcherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRatios",
			Handler:    _CurrencyFetcher_GetRatios_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ratio.proto",
}