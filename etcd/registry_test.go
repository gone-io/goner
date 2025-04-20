package etcd

import (
	"context"
	"errors"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	etcd3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	lease := NewMockLease(controller)
	kv := NewMockKV(controller)

	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	ctxClient := etcd3.NewCtxClient(timeout)
	defer cancelFunc()
	ctxClient.Lease = lease
	ctxClient.KV = kv

	service := g.NewService("test", "127.0.0.1", 8080, nil, true, 1)

	r := &Registry{
		Flag:         gone.Flag{},
		logger:       gone.GetDefaultLogger(),
		client:       ctxClient,
		dialTimeout:  10 * time.Second,
		keepaliveTTL: 2 * time.Second,
		lease:        lease,
	}

	t.Run("redoRegisterLease", func(t *testing.T) {
		lease.EXPECT().Grant(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))
		lease.EXPECT().Grant(gomock.Any(), gomock.Any()).Return(&etcd3.LeaseGrantResponse{
			ID: 1,
		}, nil)

		kv.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		lease.EXPECT().KeepAlive(gomock.Any(), gomock.Any()).Return(nil, nil)

		r.redoRegisterLease(service, etcd3.LeaseID(1))
	})

	t.Run("doRegisterLease", func(t *testing.T) {
		t.Run("put error", func(t *testing.T) {
			lease.EXPECT().Grant(gomock.Any(), gomock.Any()).Return(&etcd3.LeaseGrantResponse{
				ID: 1,
			}, nil)

			kv.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))
			err := r.doRegisterLease(context.Background(), service)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "test")
		})
		t.Run("KeepAlive error", func(t *testing.T) {
			lease.EXPECT().Grant(gomock.Any(), gomock.Any()).Return(&etcd3.LeaseGrantResponse{
				ID: 1,
			}, nil)

			kv.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			lease.EXPECT().KeepAlive(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))
			err := r.doRegisterLease(context.Background(), service)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "test")
		})
	})

}
