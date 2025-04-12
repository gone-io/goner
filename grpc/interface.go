package grpc

import "google.golang.org/grpc"

type Client interface {
	Address() string
	Stub(conn *grpc.ClientConn)
}

type Service interface {
	RegisterGrpcServer(server *grpc.Server)
}
