// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: service.proto

package service

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

// AuthenticationServiceClient is the client API for AuthenticationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthenticationServiceClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error)
}

type authenticationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthenticationServiceClient(cc grpc.ClientConnInterface) AuthenticationServiceClient {
	return &authenticationServiceClient{cc}
}

func (c *authenticationServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/go_grpc.AuthenticationService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/go_grpc.AuthenticationService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServiceClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	out := new(LogoutResponse)
	err := c.cc.Invoke(ctx, "/go_grpc.AuthenticationService/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthenticationServiceServer is the server API for AuthenticationService service.
// All implementations must embed UnimplementedAuthenticationServiceServer
// for forward compatibility
type AuthenticationServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
	mustEmbedUnimplementedAuthenticationServiceServer()
}

// UnimplementedAuthenticationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthenticationServiceServer struct {
}

func (UnimplementedAuthenticationServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthenticationServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAuthenticationServiceServer) Logout(context.Context, *LogoutRequest) (*LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedAuthenticationServiceServer) mustEmbedUnimplementedAuthenticationServiceServer() {}

// UnsafeAuthenticationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthenticationServiceServer will
// result in compilation errors.
type UnsafeAuthenticationServiceServer interface {
	mustEmbedUnimplementedAuthenticationServiceServer()
}

func RegisterAuthenticationServiceServer(s grpc.ServiceRegistrar, srv AuthenticationServiceServer) {
	s.RegisterService(&AuthenticationService_ServiceDesc, srv)
}

func _AuthenticationService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.AuthenticationService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.AuthenticationService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.AuthenticationService/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServiceServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthenticationService_ServiceDesc is the grpc.ServiceDesc for AuthenticationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthenticationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "go_grpc.AuthenticationService",
	HandlerType: (*AuthenticationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _AuthenticationService_Login_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _AuthenticationService_Register_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _AuthenticationService_Logout_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

// UsersServiceClient is the client API for UsersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersServiceClient interface {
	GetAllProfiles(ctx context.Context, in *GetAllProfilesRequest, opts ...grpc.CallOption) (*GetAllProfilesResponse, error)
	GetProfileById(ctx context.Context, in *GetProfileByIdRequest, opts ...grpc.CallOption) (*UserProfile, error)
}

type usersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersServiceClient(cc grpc.ClientConnInterface) UsersServiceClient {
	return &usersServiceClient{cc}
}

func (c *usersServiceClient) GetAllProfiles(ctx context.Context, in *GetAllProfilesRequest, opts ...grpc.CallOption) (*GetAllProfilesResponse, error) {
	out := new(GetAllProfilesResponse)
	err := c.cc.Invoke(ctx, "/go_grpc.UsersService/GetAllProfiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) GetProfileById(ctx context.Context, in *GetProfileByIdRequest, opts ...grpc.CallOption) (*UserProfile, error) {
	out := new(UserProfile)
	err := c.cc.Invoke(ctx, "/go_grpc.UsersService/GetProfileById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServiceServer is the server API for UsersService service.
// All implementations must embed UnimplementedUsersServiceServer
// for forward compatibility
type UsersServiceServer interface {
	GetAllProfiles(context.Context, *GetAllProfilesRequest) (*GetAllProfilesResponse, error)
	GetProfileById(context.Context, *GetProfileByIdRequest) (*UserProfile, error)
	mustEmbedUnimplementedUsersServiceServer()
}

// UnimplementedUsersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUsersServiceServer struct {
}

func (UnimplementedUsersServiceServer) GetAllProfiles(context.Context, *GetAllProfilesRequest) (*GetAllProfilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllProfiles not implemented")
}
func (UnimplementedUsersServiceServer) GetProfileById(context.Context, *GetProfileByIdRequest) (*UserProfile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfileById not implemented")
}
func (UnimplementedUsersServiceServer) mustEmbedUnimplementedUsersServiceServer() {}

// UnsafeUsersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersServiceServer will
// result in compilation errors.
type UnsafeUsersServiceServer interface {
	mustEmbedUnimplementedUsersServiceServer()
}

func RegisterUsersServiceServer(s grpc.ServiceRegistrar, srv UsersServiceServer) {
	s.RegisterService(&UsersService_ServiceDesc, srv)
}

func _UsersService_GetAllProfiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllProfilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetAllProfiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.UsersService/GetAllProfiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetAllProfiles(ctx, req.(*GetAllProfilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_GetProfileById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetProfileById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.UsersService/GetProfileById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetProfileById(ctx, req.(*GetProfileByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UsersService_ServiceDesc is the grpc.ServiceDesc for UsersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UsersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "go_grpc.UsersService",
	HandlerType: (*UsersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllProfiles",
			Handler:    _UsersService_GetAllProfiles_Handler,
		},
		{
			MethodName: "GetProfileById",
			Handler:    _UsersService_GetProfileById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

// ProductsServiceClient is the client API for ProductsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProductsServiceClient interface {
	CreateProduct(ctx context.Context, in *CreateProductRequest, opts ...grpc.CallOption) (*ProductResponse, error)
	GetAllProducts(ctx context.Context, in *GetAllProductsRequest, opts ...grpc.CallOption) (*GetAllProductsResponse, error)
	GetProductById(ctx context.Context, in *GetProductByIdRequest, opts ...grpc.CallOption) (*ProductResponse, error)
	GetProductsByStock(ctx context.Context, in *GetProductsByStockRequest, opts ...grpc.CallOption) (*GetProductsByStockResponse, error)
	UpdateProduct(ctx context.Context, in *UpdateProductRequest, opts ...grpc.CallOption) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, in *DeleteProductRequest, opts ...grpc.CallOption) (*DeleteProductResponse, error)
}

type productsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProductsServiceClient(cc grpc.ClientConnInterface) ProductsServiceClient {
	return &productsServiceClient{cc}
}

func (c *productsServiceClient) CreateProduct(ctx context.Context, in *CreateProductRequest, opts ...grpc.CallOption) (*ProductResponse, error) {
	out := new(ProductResponse)
	err := c.cc.Invoke(ctx, "/go_grpc.ProductsService/CreateProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productsServiceClient) GetAllProducts(ctx context.Context, in *GetAllProductsRequest, opts ...grpc.CallOption) (*GetAllProductsResponse, error) {
	out := new(GetAllProductsResponse)
	err := c.cc.Invoke(ctx, "/go_grpc.ProductsService/GetAllProducts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productsServiceClient) GetProductById(ctx context.Context, in *GetProductByIdRequest, opts ...grpc.CallOption) (*ProductResponse, error) {
	out := new(ProductResponse)
	err := c.cc.Invoke(ctx, "/go_grpc.ProductsService/GetProductById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productsServiceClient) GetProductsByStock(ctx context.Context, in *GetProductsByStockRequest, opts ...grpc.CallOption) (*GetProductsByStockResponse, error) {
	out := new(GetProductsByStockResponse)
	err := c.cc.Invoke(ctx, "/go_grpc.ProductsService/GetProductsByStock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productsServiceClient) UpdateProduct(ctx context.Context, in *UpdateProductRequest, opts ...grpc.CallOption) (*ProductResponse, error) {
	out := new(ProductResponse)
	err := c.cc.Invoke(ctx, "/go_grpc.ProductsService/UpdateProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productsServiceClient) DeleteProduct(ctx context.Context, in *DeleteProductRequest, opts ...grpc.CallOption) (*DeleteProductResponse, error) {
	out := new(DeleteProductResponse)
	err := c.cc.Invoke(ctx, "/go_grpc.ProductsService/DeleteProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductsServiceServer is the server API for ProductsService service.
// All implementations must embed UnimplementedProductsServiceServer
// for forward compatibility
type ProductsServiceServer interface {
	CreateProduct(context.Context, *CreateProductRequest) (*ProductResponse, error)
	GetAllProducts(context.Context, *GetAllProductsRequest) (*GetAllProductsResponse, error)
	GetProductById(context.Context, *GetProductByIdRequest) (*ProductResponse, error)
	GetProductsByStock(context.Context, *GetProductsByStockRequest) (*GetProductsByStockResponse, error)
	UpdateProduct(context.Context, *UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(context.Context, *DeleteProductRequest) (*DeleteProductResponse, error)
	mustEmbedUnimplementedProductsServiceServer()
}

// UnimplementedProductsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProductsServiceServer struct {
}

func (UnimplementedProductsServiceServer) CreateProduct(context.Context, *CreateProductRequest) (*ProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProduct not implemented")
}
func (UnimplementedProductsServiceServer) GetAllProducts(context.Context, *GetAllProductsRequest) (*GetAllProductsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllProducts not implemented")
}
func (UnimplementedProductsServiceServer) GetProductById(context.Context, *GetProductByIdRequest) (*ProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProductById not implemented")
}
func (UnimplementedProductsServiceServer) GetProductsByStock(context.Context, *GetProductsByStockRequest) (*GetProductsByStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProductsByStock not implemented")
}
func (UnimplementedProductsServiceServer) UpdateProduct(context.Context, *UpdateProductRequest) (*ProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProduct not implemented")
}
func (UnimplementedProductsServiceServer) DeleteProduct(context.Context, *DeleteProductRequest) (*DeleteProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProduct not implemented")
}
func (UnimplementedProductsServiceServer) mustEmbedUnimplementedProductsServiceServer() {}

// UnsafeProductsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProductsServiceServer will
// result in compilation errors.
type UnsafeProductsServiceServer interface {
	mustEmbedUnimplementedProductsServiceServer()
}

func RegisterProductsServiceServer(s grpc.ServiceRegistrar, srv ProductsServiceServer) {
	s.RegisterService(&ProductsService_ServiceDesc, srv)
}

func _ProductsService_CreateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).CreateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.ProductsService/CreateProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).CreateProduct(ctx, req.(*CreateProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductsService_GetAllProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllProductsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).GetAllProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.ProductsService/GetAllProducts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).GetAllProducts(ctx, req.(*GetAllProductsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductsService_GetProductById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).GetProductById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.ProductsService/GetProductById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).GetProductById(ctx, req.(*GetProductByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductsService_GetProductsByStock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductsByStockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).GetProductsByStock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.ProductsService/GetProductsByStock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).GetProductsByStock(ctx, req.(*GetProductsByStockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductsService_UpdateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).UpdateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.ProductsService/UpdateProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).UpdateProduct(ctx, req.(*UpdateProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductsService_DeleteProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).DeleteProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_grpc.ProductsService/DeleteProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).DeleteProduct(ctx, req.(*DeleteProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProductsService_ServiceDesc is the grpc.ServiceDesc for ProductsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProductsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "go_grpc.ProductsService",
	HandlerType: (*ProductsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProduct",
			Handler:    _ProductsService_CreateProduct_Handler,
		},
		{
			MethodName: "GetAllProducts",
			Handler:    _ProductsService_GetAllProducts_Handler,
		},
		{
			MethodName: "GetProductById",
			Handler:    _ProductsService_GetProductById_Handler,
		},
		{
			MethodName: "GetProductsByStock",
			Handler:    _ProductsService_GetProductsByStock_Handler,
		},
		{
			MethodName: "UpdateProduct",
			Handler:    _ProductsService_UpdateProduct_Handler,
		},
		{
			MethodName: "DeleteProduct",
			Handler:    _ProductsService_DeleteProduct_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
