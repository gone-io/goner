package etcd

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	etcd3 "go.etcd.io/etcd/client/v3"
)

var client *etcd3.Client

func ProvideEtecd3Client(_ string, param struct {
	config *etcd3.Config `gone:"config,etcd"`
}) (*etcd3.Client, error) {
	if client != nil {
		return client, nil
	}
	var err error
	client, err = etcd3.New(*param.config)
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, "can not create etcd client")
	}
	return client, nil
}
