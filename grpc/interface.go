package grpc

import "google.golang.org/grpc"

//go:generate sh -c "mockgen -package=grpc net Listener > net_listener_mock_test.go"

//go:generate sh -c "mockgen -package=grpc -self_package=github.com/gone-io/goner/grpc -source=interface.go -destination=mock_test.go"

//go:generate mockgen -package=grpc -destination ./tracer_mock_test.go github.com/gone-io/goner/g Tracer,Cmux

type Client interface {
	Address() string
	Stub(conn *grpc.ClientConn)
}

type Service interface {
	RegisterGrpcServer(server *grpc.Server)
}
