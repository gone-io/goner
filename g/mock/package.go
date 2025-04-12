package mock

//go:generate mockgen -package=mock -source=../balancer.go -destination=./balancer_mock.go
//go:generate mockgen -package=mock -source=../cmux.go -destination=./cmux_mock.go
//go:generate mockgen -package=mock -source=../discovery.go -destination=./discovery_mock.go
//go:generate mockgen -package=mock -source=../locker.go -destination=./locker_mock.go
//go:generate mockgen -package=mock -source=../registry.go -destination=./registry_mock.go
//go:generate mockgen -package=mock -source=../service.go -destination=./service_mock.go
//go:generate mockgen -package=mock -source=../tracer.go -destination=./tracer_mock.go
