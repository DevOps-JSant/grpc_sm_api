// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.31.1
// source: students.proto

package grpcapipb

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
	StudentService_GetStudents_FullMethodName    = "/main.StudentService/GetStudents"
	StudentService_AddStudents_FullMethodName    = "/main.StudentService/AddStudents"
	StudentService_UpdateStudents_FullMethodName = "/main.StudentService/UpdateStudents"
	StudentService_DeleteStudents_FullMethodName = "/main.StudentService/DeleteStudents"
)

// StudentServiceClient is the client API for StudentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StudentServiceClient interface {
	GetStudents(ctx context.Context, in *GetStudentsRequest, opts ...grpc.CallOption) (*Students, error)
	AddStudents(ctx context.Context, in *Students, opts ...grpc.CallOption) (*Students, error)
	UpdateStudents(ctx context.Context, in *Students, opts ...grpc.CallOption) (*Students, error)
	DeleteStudents(ctx context.Context, in *StudentIds, opts ...grpc.CallOption) (*DeleteStudentsConfirmation, error)
}

type studentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStudentServiceClient(cc grpc.ClientConnInterface) StudentServiceClient {
	return &studentServiceClient{cc}
}

func (c *studentServiceClient) GetStudents(ctx context.Context, in *GetStudentsRequest, opts ...grpc.CallOption) (*Students, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Students)
	err := c.cc.Invoke(ctx, StudentService_GetStudents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentServiceClient) AddStudents(ctx context.Context, in *Students, opts ...grpc.CallOption) (*Students, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Students)
	err := c.cc.Invoke(ctx, StudentService_AddStudents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentServiceClient) UpdateStudents(ctx context.Context, in *Students, opts ...grpc.CallOption) (*Students, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Students)
	err := c.cc.Invoke(ctx, StudentService_UpdateStudents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentServiceClient) DeleteStudents(ctx context.Context, in *StudentIds, opts ...grpc.CallOption) (*DeleteStudentsConfirmation, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteStudentsConfirmation)
	err := c.cc.Invoke(ctx, StudentService_DeleteStudents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StudentServiceServer is the server API for StudentService service.
// All implementations must embed UnimplementedStudentServiceServer
// for forward compatibility.
type StudentServiceServer interface {
	GetStudents(context.Context, *GetStudentsRequest) (*Students, error)
	AddStudents(context.Context, *Students) (*Students, error)
	UpdateStudents(context.Context, *Students) (*Students, error)
	DeleteStudents(context.Context, *StudentIds) (*DeleteStudentsConfirmation, error)
	mustEmbedUnimplementedStudentServiceServer()
}

// UnimplementedStudentServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStudentServiceServer struct{}

func (UnimplementedStudentServiceServer) GetStudents(context.Context, *GetStudentsRequest) (*Students, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudents not implemented")
}
func (UnimplementedStudentServiceServer) AddStudents(context.Context, *Students) (*Students, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddStudents not implemented")
}
func (UnimplementedStudentServiceServer) UpdateStudents(context.Context, *Students) (*Students, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStudents not implemented")
}
func (UnimplementedStudentServiceServer) DeleteStudents(context.Context, *StudentIds) (*DeleteStudentsConfirmation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStudents not implemented")
}
func (UnimplementedStudentServiceServer) mustEmbedUnimplementedStudentServiceServer() {}
func (UnimplementedStudentServiceServer) testEmbeddedByValue()                        {}

// UnsafeStudentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StudentServiceServer will
// result in compilation errors.
type UnsafeStudentServiceServer interface {
	mustEmbedUnimplementedStudentServiceServer()
}

func RegisterStudentServiceServer(s grpc.ServiceRegistrar, srv StudentServiceServer) {
	// If the following call pancis, it indicates UnimplementedStudentServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StudentService_ServiceDesc, srv)
}

func _StudentService_GetStudents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStudentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).GetStudents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudentService_GetStudents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).GetStudents(ctx, req.(*GetStudentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudentService_AddStudents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Students)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).AddStudents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudentService_AddStudents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).AddStudents(ctx, req.(*Students))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudentService_UpdateStudents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Students)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).UpdateStudents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudentService_UpdateStudents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).UpdateStudents(ctx, req.(*Students))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudentService_DeleteStudents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StudentIds)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).DeleteStudents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudentService_DeleteStudents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).DeleteStudents(ctx, req.(*StudentIds))
	}
	return interceptor(ctx, in, info, handler)
}

// StudentService_ServiceDesc is the grpc.ServiceDesc for StudentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StudentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.StudentService",
	HandlerType: (*StudentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStudents",
			Handler:    _StudentService_GetStudents_Handler,
		},
		{
			MethodName: "AddStudents",
			Handler:    _StudentService_AddStudents_Handler,
		},
		{
			MethodName: "UpdateStudents",
			Handler:    _StudentService_UpdateStudents_Handler,
		},
		{
			MethodName: "DeleteStudents",
			Handler:    _StudentService_DeleteStudents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "students.proto",
}
