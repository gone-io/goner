package xorm

import (
	"errors"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_trans_Transaction(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := gone.GetDefaultLogger()

	session := NewMockSession(controller)

	tests := []struct {
		name       string
		newSession func() Session
		fn         func(session trans) error
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			newSession: func() Session {
				session.EXPECT().Begin().Return(nil)
				session.EXPECT().Exec(`select * from test`).Return(nil, nil)
				session.EXPECT().Commit().Return(nil)
				session.EXPECT().Close().Return(nil)
				return session
			},
			fn: func(x trans) error {
				return x.Transaction(func(session Interface) error {
					exec, err := session.Exec(`select * from test`)
					assert.Nil(t, err)
					assert.Nil(t, exec)
					return nil
				})
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "nest transaction success",
			newSession: func() Session {
				session.EXPECT().Begin().Return(nil)
				session.EXPECT().Exec(`select * from test`).Return(nil, nil)
				session.EXPECT().Exec(`select * from test2`).Return(nil, nil)
				session.EXPECT().Commit().Return(nil)
				session.EXPECT().Close().Return(nil)
				return session
			},
			fn: func(x trans) error {
				return x.Transaction(func(session Interface) error {
					exec, err := session.Exec(`select * from test`)
					assert.Nil(t, err)
					assert.Nil(t, exec)

					return x.Transaction(func(session Interface) error {
						_, _ = session.Exec(`select * from test2`)
						return nil
					})
				})
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "err in begin",
			newSession: func() Session {
				session.EXPECT().Begin().Return(errors.New("begin err"))
				session.EXPECT().Close().Return(nil)
				return session
			},
			fn: func(x trans) error {
				return x.Transaction(func(session Interface) error {
					return nil
				})
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				var goneErr gone.Error
				assert.True(t, errors.As(err, &goneErr))
				return assert.Contains(t, goneErr.Error(), "begin err")
			},
		},
		{
			name: "err in process fn",
			newSession: func() Session {
				session.EXPECT().Begin().Return(nil)
				session.EXPECT().Rollback().Return(nil)
				session.EXPECT().Close().Return(nil)
				return session
			},
			fn: func(x trans) error {
				return x.Transaction(func(session Interface) error {
					return errors.New("process err")
				})
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				var goneErr gone.Error
				assert.True(t, errors.As(err, &goneErr))
				return assert.Contains(t, goneErr.Error(), "process err")
			},
		},
		{
			name: "err in process fn in nest transaction",
			newSession: func() Session {
				session.EXPECT().Begin().Return(nil)
				session.EXPECT().Rollback().Return(nil)
				session.EXPECT().Close().Return(nil)
				return session
			},
			fn: func(x trans) error {
				return x.Transaction(func(session Interface) error {
					return x.Transaction(func(session Interface) error {
						return errors.New("process err")
					})
				})
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				var goneErr gone.Error
				assert.True(t, errors.As(err, &goneErr))
				return assert.Contains(t, goneErr.Error(), "process err")
			},
		},
		{
			name: "err in rollback",
			newSession: func() Session {
				session.EXPECT().Begin().Return(nil)
				session.EXPECT().Rollback().Return(errors.New("test error"))
				session.EXPECT().Close().Return(nil)
				return session
			},
			fn: func(x trans) error {
				return x.Transaction(func(session Interface) error {
					return errors.New("process err")
				})
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				var goneErr gone.Error
				assert.True(t, errors.As(err, &goneErr))
				assert.Contains(t, goneErr.Error(), "test error")
				return assert.Contains(t, goneErr.Error(), "process err")
			},
		},
		{
			name: "panic in process",
			newSession: func() Session {
				session.EXPECT().Begin().Return(nil)
				session.EXPECT().Rollback().Return(nil)
				session.EXPECT().Close().Return(nil)
				return session
			},
			fn: func(x trans) error {
				return x.Transaction(func(session Interface) error {
					panic("test panic")
				})
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				var goneErr gone.Error
				assert.True(t, errors.As(err, &goneErr))
				return assert.Contains(t, goneErr.Error(), "test panic")
			},
		},
		{
			name: "panic in rollback",
			newSession: func() Session {
				session.EXPECT().Begin().Return(nil)
				session.EXPECT().Rollback().Do(func() {
					panic("test panic")
				})
				session.EXPECT().Close().Return(nil)
				return session
			},
			fn: func(x trans) error {
				return x.Transaction(func(session Interface) error {
					return errors.New("test error")
				})
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				var goneErr gone.Error
				assert.True(t, errors.As(err, &goneErr))
				assert.Contains(t, goneErr.Error(), "test error")
				return assert.Contains(t, goneErr.Error(), "test panic")
			},
		},
		{
			name: "err when close",
			newSession: func() Session {
				session.EXPECT().Begin().Return(nil)
				session.EXPECT().Commit().Return(nil)
				session.EXPECT().Close().Return(errors.New("test error"))
				return session
			},
			fn: func(session trans) error {
				return session.Transaction(func(session Interface) error {
					return nil
				})
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newTrans(logger, tt.newSession)
			tt.wantErr(t, tt.fn(e), fmt.Sprintf("Transaction(%T)", tt.fn))
		})
	}
}
