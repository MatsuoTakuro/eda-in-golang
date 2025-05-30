// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: customerspb/api.proto

package customerspb

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
	CustomersService_RegisterCustomer_FullMethodName  = "/customerspb.CustomersService/RegisterCustomer"
	CustomersService_EnableCustomer_FullMethodName    = "/customerspb.CustomersService/EnableCustomer"
	CustomersService_DisableCustomer_FullMethodName   = "/customerspb.CustomersService/DisableCustomer"
	CustomersService_ChangeSmsNumber_FullMethodName   = "/customerspb.CustomersService/ChangeSmsNumber"
	CustomersService_AuthorizeCustomer_FullMethodName = "/customerspb.CustomersService/AuthorizeCustomer"
	CustomersService_GetCustomer_FullMethodName       = "/customerspb.CustomersService/GetCustomer"
)

// CustomersServiceClient is the client API for CustomersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CustomersServiceClient interface {
	RegisterCustomer(ctx context.Context, in *RegisterCustomerRequest, opts ...grpc.CallOption) (*RegisterCustomerResponse, error)
	EnableCustomer(ctx context.Context, in *EnableCustomerRequest, opts ...grpc.CallOption) (*EnableCustomerResponse, error)
	DisableCustomer(ctx context.Context, in *DisableCustomerRequest, opts ...grpc.CallOption) (*DisableCustomerResponse, error)
	ChangeSmsNumber(ctx context.Context, in *ChangeSmsNumberRequest, opts ...grpc.CallOption) (*ChangeSmsNumberResponse, error)
	AuthorizeCustomer(ctx context.Context, in *AuthorizeCustomerRequest, opts ...grpc.CallOption) (*AuthorizeCustomerResponse, error)
	GetCustomer(ctx context.Context, in *GetCustomerRequest, opts ...grpc.CallOption) (*GetCustomerResponse, error)
}

type customersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCustomersServiceClient(cc grpc.ClientConnInterface) CustomersServiceClient {
	return &customersServiceClient{cc}
}

func (c *customersServiceClient) RegisterCustomer(ctx context.Context, in *RegisterCustomerRequest, opts ...grpc.CallOption) (*RegisterCustomerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterCustomerResponse)
	err := c.cc.Invoke(ctx, CustomersService_RegisterCustomer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customersServiceClient) EnableCustomer(ctx context.Context, in *EnableCustomerRequest, opts ...grpc.CallOption) (*EnableCustomerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EnableCustomerResponse)
	err := c.cc.Invoke(ctx, CustomersService_EnableCustomer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customersServiceClient) DisableCustomer(ctx context.Context, in *DisableCustomerRequest, opts ...grpc.CallOption) (*DisableCustomerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DisableCustomerResponse)
	err := c.cc.Invoke(ctx, CustomersService_DisableCustomer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customersServiceClient) ChangeSmsNumber(ctx context.Context, in *ChangeSmsNumberRequest, opts ...grpc.CallOption) (*ChangeSmsNumberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ChangeSmsNumberResponse)
	err := c.cc.Invoke(ctx, CustomersService_ChangeSmsNumber_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customersServiceClient) AuthorizeCustomer(ctx context.Context, in *AuthorizeCustomerRequest, opts ...grpc.CallOption) (*AuthorizeCustomerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthorizeCustomerResponse)
	err := c.cc.Invoke(ctx, CustomersService_AuthorizeCustomer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customersServiceClient) GetCustomer(ctx context.Context, in *GetCustomerRequest, opts ...grpc.CallOption) (*GetCustomerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCustomerResponse)
	err := c.cc.Invoke(ctx, CustomersService_GetCustomer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CustomersServiceServer is the server API for CustomersService service.
// All implementations must embed UnimplementedCustomersServiceServer
// for forward compatibility.
type CustomersServiceServer interface {
	RegisterCustomer(context.Context, *RegisterCustomerRequest) (*RegisterCustomerResponse, error)
	EnableCustomer(context.Context, *EnableCustomerRequest) (*EnableCustomerResponse, error)
	DisableCustomer(context.Context, *DisableCustomerRequest) (*DisableCustomerResponse, error)
	ChangeSmsNumber(context.Context, *ChangeSmsNumberRequest) (*ChangeSmsNumberResponse, error)
	AuthorizeCustomer(context.Context, *AuthorizeCustomerRequest) (*AuthorizeCustomerResponse, error)
	GetCustomer(context.Context, *GetCustomerRequest) (*GetCustomerResponse, error)
	mustEmbedUnimplementedCustomersServiceServer()
}

// UnimplementedCustomersServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCustomersServiceServer struct{}

func (UnimplementedCustomersServiceServer) RegisterCustomer(context.Context, *RegisterCustomerRequest) (*RegisterCustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterCustomer not implemented")
}
func (UnimplementedCustomersServiceServer) EnableCustomer(context.Context, *EnableCustomerRequest) (*EnableCustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableCustomer not implemented")
}
func (UnimplementedCustomersServiceServer) DisableCustomer(context.Context, *DisableCustomerRequest) (*DisableCustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableCustomer not implemented")
}
func (UnimplementedCustomersServiceServer) ChangeSmsNumber(context.Context, *ChangeSmsNumberRequest) (*ChangeSmsNumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeSmsNumber not implemented")
}
func (UnimplementedCustomersServiceServer) AuthorizeCustomer(context.Context, *AuthorizeCustomerRequest) (*AuthorizeCustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthorizeCustomer not implemented")
}
func (UnimplementedCustomersServiceServer) GetCustomer(context.Context, *GetCustomerRequest) (*GetCustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCustomer not implemented")
}
func (UnimplementedCustomersServiceServer) mustEmbedUnimplementedCustomersServiceServer() {}
func (UnimplementedCustomersServiceServer) testEmbeddedByValue()                          {}

// UnsafeCustomersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CustomersServiceServer will
// result in compilation errors.
type UnsafeCustomersServiceServer interface {
	mustEmbedUnimplementedCustomersServiceServer()
}

func RegisterCustomersServiceServer(s grpc.ServiceRegistrar, srv CustomersServiceServer) {
	// If the following call pancis, it indicates UnimplementedCustomersServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CustomersService_ServiceDesc, srv)
}

func _CustomersService_RegisterCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).RegisterCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CustomersService_RegisterCustomer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).RegisterCustomer(ctx, req.(*RegisterCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomersService_EnableCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnableCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).EnableCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CustomersService_EnableCustomer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).EnableCustomer(ctx, req.(*EnableCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomersService_DisableCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisableCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).DisableCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CustomersService_DisableCustomer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).DisableCustomer(ctx, req.(*DisableCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomersService_ChangeSmsNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeSmsNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).ChangeSmsNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CustomersService_ChangeSmsNumber_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).ChangeSmsNumber(ctx, req.(*ChangeSmsNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomersService_AuthorizeCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizeCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).AuthorizeCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CustomersService_AuthorizeCustomer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).AuthorizeCustomer(ctx, req.(*AuthorizeCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomersService_GetCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).GetCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CustomersService_GetCustomer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).GetCustomer(ctx, req.(*GetCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CustomersService_ServiceDesc is the grpc.ServiceDesc for CustomersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CustomersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "customerspb.CustomersService",
	HandlerType: (*CustomersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterCustomer",
			Handler:    _CustomersService_RegisterCustomer_Handler,
		},
		{
			MethodName: "EnableCustomer",
			Handler:    _CustomersService_EnableCustomer_Handler,
		},
		{
			MethodName: "DisableCustomer",
			Handler:    _CustomersService_DisableCustomer_Handler,
		},
		{
			MethodName: "ChangeSmsNumber",
			Handler:    _CustomersService_ChangeSmsNumber_Handler,
		},
		{
			MethodName: "AuthorizeCustomer",
			Handler:    _CustomersService_AuthorizeCustomer_Handler,
		},
		{
			MethodName: "GetCustomer",
			Handler:    _CustomersService_GetCustomer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "customerspb/api.proto",
}
