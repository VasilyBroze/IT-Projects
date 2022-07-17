// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.2
// source: api/proto/finder.proto

package gRPC_Books_Test

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

// BooksAndAuthorsClient is the client API for BooksAndAuthors service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BooksAndAuthorsClient interface {
	GetAuthor(ctx context.Context, in *AuthorRequest, opts ...grpc.CallOption) (*AuthorResponse, error)
}

type booksAndAuthorsClient struct {
	cc grpc.ClientConnInterface
}

func NewBooksAndAuthorsClient(cc grpc.ClientConnInterface) BooksAndAuthorsClient {
	return &booksAndAuthorsClient{cc}
}

func (c *booksAndAuthorsClient) GetAuthor(ctx context.Context, in *AuthorRequest, opts ...grpc.CallOption) (*AuthorResponse, error) {
	out := new(AuthorResponse)
	err := c.cc.Invoke(ctx, "/api.BooksAndAuthors/GetAuthor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BooksAndAuthorsServer is the server API for BooksAndAuthors service.
// All implementations must embed UnimplementedBooksAndAuthorsServer
// for forward compatibility
type BooksAndAuthorsServer interface {
	GetAuthor(context.Context, *AuthorRequest) (*AuthorResponse, error)
	mustEmbedUnimplementedBooksAndAuthorsServer()
}

// UnimplementedBooksAndAuthorsServer must be embedded to have forward compatible implementations.
type UnimplementedBooksAndAuthorsServer struct {
}

func (UnimplementedBooksAndAuthorsServer) GetAuthor(context.Context, *AuthorRequest) (*AuthorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthor not implemented")
}
func (UnimplementedBooksAndAuthorsServer) mustEmbedUnimplementedBooksAndAuthorsServer() {}

// UnsafeBooksAndAuthorsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BooksAndAuthorsServer will
// result in compilation errors.
type UnsafeBooksAndAuthorsServer interface {
	mustEmbedUnimplementedBooksAndAuthorsServer()
}

func RegisterBooksAndAuthorsServer(s grpc.ServiceRegistrar, srv BooksAndAuthorsServer) {
	s.RegisterService(&BooksAndAuthors_ServiceDesc, srv)
}

func _BooksAndAuthors_GetAuthor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BooksAndAuthorsServer).GetAuthor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.BooksAndAuthors/GetAuthor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BooksAndAuthorsServer).GetAuthor(ctx, req.(*AuthorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BooksAndAuthors_ServiceDesc is the grpc.ServiceDesc for BooksAndAuthors service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BooksAndAuthors_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.BooksAndAuthors",
	HandlerType: (*BooksAndAuthorsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAuthor",
			Handler:    _BooksAndAuthors_GetAuthor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/finder.proto",
}