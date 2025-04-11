package grpc

//go:generate mockgen -package=grpc -destination=./net_mock.go net Listener,Addr,Conn

//go:generate mockgen -package=grpc -self_package=github.com/gone-io/goner/grpc -source=interface.go -destination=grpc_mock.go

//go:generate mockgen -package=grpc -destination=grpc_resolver_mock.go google.golang.org/grpc/resolver ClientConn
