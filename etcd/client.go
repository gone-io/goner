package etcd

import (
	"github.com/gone-io/gone/v2"
	etcd3 "go.etcd.io/etcd/client/v3"
)

var client *etcd3.Client

func ProvideEtecd3Client(_ string, param struct {
	config *etcd3.Config `gone:"config,etcd"`
	conf   *etcd3.Config `gone:"etcd.config" option:"allowNil"`
}) (*etcd3.Client, error) {
	if client != nil {
		return client, nil
	}
	var err error
	if param.conf != nil {
		param.config = param.conf
	}
	client, err = etcd3.New(*param.config)
	return client, gone.ToErrorWithMsg(err, "can not create etcd client")
}
