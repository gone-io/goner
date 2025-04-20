package etcd

import (
	"context"
	"errors"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/api/v3/mvccpb"
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

	t.Run("GetInstances Get error", func(t *testing.T) {
		kv.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))
		_, err := r.GetInstances("test")
		assert.Error(t, err)
	})

	t.Run("GetInstances extractResponseToServices err", func(t *testing.T) {
		kv.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(&etcd3.GetResponse{
			Kvs: []*mvccpb.KeyValue{
				{
					Key:   []byte("test"),
					Value: []byte(``),
				},
			},
		}, nil)
		_, err := r.GetInstances("test")
		assert.Error(t, err)
	})

	t.Run("Watch connect error", func(t *testing.T) {
		kv.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))
		_, _, err := r.Watch("test")
		assert.Error(t, err)
	})

	t.Run("watch error", func(t *testing.T) {
		watcher := NewMockWatcher(controller)
		var ch etcd3.WatchChan
		watcher.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).Return(ch)
		watcher.EXPECT().RequestProgress(gomock.Any()).Return(errors.New("test"))
		_, _, err := r.watch("test", watcher)
		assert.Error(t, err)
	})

	t.Run("watch success", func(t *testing.T) {
		watcher := NewMockWatcher(controller)
		var ch = make(chan etcd3.WatchResponse)
		watcher.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).Return(ch)
		watcher.EXPECT().RequestProgress(gomock.Any()).Return(nil)
		watcher.EXPECT().Close().Return(nil)
		kv.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(&etcd3.GetResponse{
			Kvs: []*mvccpb.KeyValue{
				{
					Key:   []byte("test"),
					Value: []byte(g.GetServerValue(g.NewService("test", "127.0.0.1", 8080, nil, true, 1))),
				},
			},
		}, nil)

		sCh, stop, err := r.watch("test", watcher)
		assert.Nil(t, err)
		defer stop()
		assert.NotNil(t, sCh)
		go func() {
			ch <- etcd3.WatchResponse{}
		}()
		services := <-sCh
		assert.Len(t, services, 1)
		assert.Equal(t, "test", services[0].GetName())

	})

}
