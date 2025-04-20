package etcd

import (
	"fmt"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd3 "go.etcd.io/etcd/client/v3"
	"testing"
)

func Test_extractResponseToServices(t *testing.T) {
	type args struct {
		res *etcd3.GetResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "suc",
			args: args{
				res: &etcd3.GetResponse{
					Kvs: []*mvccpb.KeyValue{
						{
							Key:   []byte("test"),
							Value: []byte(g.GetServerValue(g.NewService("test", "127.0.0.1", 8080, nil, true, 1))),
						},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "failed",
			args: args{
				res: &etcd3.GetResponse{
					Kvs: []*mvccpb.KeyValue{
						{
							Key:   []byte("test"),
							Value: []byte(`{"name":"test"`),
						},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := extractResponseToServices(tt.args.res)
			if !tt.wantErr(t, err, fmt.Sprintf("extractResponseToServices(%v)", tt.args.res)) {
				return
			}
		})
	}
}
