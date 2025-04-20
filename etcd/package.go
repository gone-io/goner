package etcd

//go:generate mockgen -destination=etcd_mock.go -package=etcd go.etcd.io/etcd/client/v3 KV,Lease
