package xorm

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"xorm.io/xorm"
)

func Test_engine(t *testing.T) {
	gone.
		Test(func(in struct {
			logger gone.Logger `gone:"gone-logger"`
		}) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			engineInterface := NewMockEngineInterface(controller)
			engineInterface.EXPECT().SetConnMaxLifetime(gomock.Any())
			engineInterface.EXPECT().SetMaxOpenConns(gomock.Any())
			engineInterface.EXPECT().SetMaxIdleConns(gomock.Any())
			engineInterface.EXPECT().SetLogger(gomock.Any())
			engineInterface.EXPECT().Ping()
			engineInterface.EXPECT().SQL(gomock.Any(), gomock.Any()).Return(nil)

			e := wrappedEngine{
				log: in.logger,
				newFunc: func(driverName string, dataSourceName string) (xorm.EngineInterface, error) {
					return nil, errors.New("test")
				},
			}

			err := e.Start()
			assert.Error(t, err)

			e.newFunc = func(driverName string, dataSourceName string) (xorm.EngineInterface, error) {
				return engineInterface, nil
			}

			err = e.Start()
			assert.NoError(t, err)

			originEngine := e.GetOriginEngine()
			assert.Equalf(t, engineInterface, originEngine, "origin wrappedEngine is not equal")

			_ = e.Sqlx("select * from user where id = ?", 1)

			err = e.Start()
			assert.Error(t, err)

		})
}

func Test_engineCluster(t *testing.T) {
	gone.Test(func(in struct {
		logger gone.Logger `gone:"gone-logger"`
	}) {
		db, mock, _ := sqlmock.New(
			sqlmock.MonitorPingsOption(true),
		)
		defer db.Close()

		sql.Register("mysql", db.Driver())

		mock.ExpectPing()
		mock.ExpectPing()
		mock.ExpectClose()

		e := wrappedEngine{
			log:           in.logger,
			enableCluster: true,
			masterConf: &ClusterNodeConf{
				DriverName: "mysql",
				DSN:        "sqlmock_db_0",
			},
			slavesConf: []*ClusterNodeConf{
				{
					DriverName: "mysql",
					DSN:        "sqlmock_db_0",
				},
			},
		}

		// 测试集群配置错误场景
		err := e.Start()
		assert.Error(t, err, "should fail when newFunc is nil")

		// 设置mock newFunc
		e.newFunc = func(driverName string, dataSourceName string) (xorm.EngineInterface, error) {
			return xorm.NewEngine(driverName, dataSourceName)
		}

		err = e.Start()
		assert.NoError(t, err)
	})
}

func Test_engineStop(t *testing.T) {
	gone.Test(func(in struct {
		logger gone.Logger `gone:"gone-logger"`
	}) {
		db, mock, _ := sqlmock.New(
			sqlmock.MonitorPingsOption(true),
		)
		defer db.Close()

		sql.Register("sqlite", db.Driver())

		mock.ExpectPing()
		mock.ExpectPing()
		mock.ExpectClose()

		e := wrappedEngine{
			log: in.logger,
			conf: Conf{
				DriverName: "sqlite",
				Dsn:        "sqlmock_db_1",
			},
			newFunc: newEngine,
		}

		err := e.Start()
		assert.NoError(t, err)

		err = e.Stop()
		assert.NoError(t, err)
	})
}

func Test_engineConfig(t *testing.T) {
	gone.Test(func(in struct {
		logger gone.Logger `gone:"gone-logger"`
	}) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		engineInterface := NewMockEngineInterface(controller)

		conf := Conf{
			DriverName:   "mysql",
			Dsn:          "test-dsn",
			MaxIdleCount: 10,
			MaxOpen:      50,
			MaxLifetime:  5 * time.Minute,
			ShowSql:      true,
		}

		e := wrappedEngine{
			log:  in.logger,
			conf: conf,
			newFunc: func(driverName string, dataSourceName string) (xorm.EngineInterface, error) {
				return engineInterface, nil
			},
		}

		engineInterface.EXPECT().SetConnMaxLifetime(conf.MaxLifetime)
		engineInterface.EXPECT().SetMaxOpenConns(conf.MaxOpen)
		engineInterface.EXPECT().SetMaxIdleConns(conf.MaxIdleCount)
		engineInterface.EXPECT().SetLogger(gomock.Any())
		engineInterface.EXPECT().Ping()

		err := e.Start()
		assert.NoError(t, err)
	})
}

func Test_engineClusterErrors(t *testing.T) {
	gone.Test(func(in struct {
		logger gone.Logger `gone:"gone-logger"`
	}) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		// 测试没有master配置的情况
		e := wrappedEngine{
			log:           in.logger,
			enableCluster: true,
		}
		err := e.Start()
		assert.Error(t, err)

		// 测试没有slave配置的情况
		e.masterConf = &ClusterNodeConf{
			DriverName: "mysql",
			DSN:        "master-dsn",
		}
		err = e.Start()
		assert.Error(t, err)

		// 测试master创建失败的情况
		e.slavesConf = []*ClusterNodeConf{{
			DriverName: "mysql",
			DSN:        "slave-dsn",
		}}
		e.newFunc = func(driverName string, dataSourceName string) (xorm.EngineInterface, error) {
			return nil, errors.New("create engine failed")
		}
		err = e.Start()
		assert.Error(t, err)
	})
}
