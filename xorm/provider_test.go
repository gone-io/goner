package xorm

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mock "github.com/gone-io/gone/mock/v2"
	"github.com/gone-io/gone/v2"
	"go.uber.org/mock/gomock"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"xorm.io/xorm"
)

type mockEngine struct {
	xorm.EngineInterface
}

func TestNewProvider(t *testing.T) {
	engine := &wrappedEngine{
		EngineInterface: &mockEngine{},
	}
	p := newProvider(engine).(*provider)

	assert.NotNil(t, p)
	assert.NotNil(t, p.engineMap)
	assert.Equal(t, engine, p.engineMap[""])
	assert.Equal(t, engine, p.engineMap[defaultCluster])
}

func TestProvider_Provide(t *testing.T) {
	engine := &wrappedEngine{
		EngineInterface: &mockEngine{},
	}
	p := newProvider(engine).(*provider)

	// Test providing XormEngine
	value, err := p.Provide("", xormInterface)
	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.IsType(t, &wrappedEngine{}, value)

	// Test providing []XormEngine with cluster disabled
	_, err = p.Provide("", xormInterfaceSlice)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not enable cluster")

	// Test providing master with cluster disabled
	_, err = p.Provide("master=true", xormInterface)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not enable cluster")

	// Test providing slave with cluster disabled
	_, err = p.Provide("slave=0", xormInterface)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not enable cluster")
}

func TestProvider_GetDb(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	configure := mock.NewMockConfigure(controller)
	configure.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to get config for cluster"))

	engine := &wrappedEngine{
		EngineInterface: &mockEngine{},
	}
	p := newProvider(engine).(*provider)
	p.configure = configure

	// Test getting existing database
	db, err := p.getDb("")
	assert.NoError(t, err)
	assert.Equal(t, engine, db)

	// Test getting non-existing database
	_, err = p.getDb("non_existing")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get config for cluster")
}

func TestProvider_GonerName(t *testing.T) {
	engine := &wrappedEngine{
		EngineInterface: &mockEngine{},
	}
	p := newProvider(engine).(*provider)

	assert.Equal(t, "xorm", p.GonerName())
}

func Test_provider_Provide(t *testing.T) {
	_ = os.Setenv("GONE_DATABASE_CLUSTER_ENABLE", "true")
	_ = os.Setenv("GONE_DATABASE_CLUSTER_MASTER", "{\"driver-name\":\"sqlite3\",\"dsn\":\"provider-db-test\"}")
	_ = os.Setenv("GONE_DATABASE_CLUSTER_SLAVES", "[{\"driver-name\":\"sqlite3\",\"dsn\":\"provider-db-test\"},{\"driver-name\":\"sqlite3\",\"dsn\":\"provider-db-test\"}]")

	_ = os.Setenv("GONE_CUSTOM_CLUSTER_ENABLE", "true")
	_ = os.Setenv("GONE_CUSTOM_CLUSTER_MASTER", "{\"driver-name\":\"sqlite3\",\"dsn\":\"provider-db-test\"}")
	_ = os.Setenv("GONE_CUSTOM_CLUSTER_SLAVES", "[{\"driver-name\":\"sqlite3\",\"dsn\":\"provider-db-test\"},{\"driver-name\":\"sqlite3\",\"dsn\":\"provider-db-test\"}]")

	defer func() {
		_ = os.Unsetenv("GONE_DATABASE_CLUSTER_ENABLE")
		_ = os.Unsetenv("GONE_DATABASE_CLUSTER_MASTER")
		_ = os.Unsetenv("GONE_DATABASE_CLUSTER_SLAVE")
		_ = os.Unsetenv("GONE_CUSTOM_CLUSTER_ENABLE")
		_ = os.Unsetenv("GONE_CUSTOM_CLUSTER_MASTER")
		_ = os.Unsetenv("GONE_CUSTOM_CLUSTER_SLAVE")
	}()

	db, dbMock, _ := sqlmock.NewWithDSN(
		"provider-db-test",
		sqlmock.MonitorPingsOption(true),
	)
	defer db.Close()

	sql.Register("sqlite3", db.Driver())

	dbMock.ExpectPing()
	dbMock.ExpectPing()
	dbMock.ExpectPing()
	dbMock.ExpectPing()
	dbMock.ExpectPing()
	dbMock.ExpectPing()

	gone.
		NewApp(Load).
		Run(func(in struct {
			db       Engine   `gone:"*"`         //to get default cluster db
			dbMaster Engine   `gone:"*,master"`  //use provider to get default cluster master db
			dbSlaves []Engine `gone:"*,slave"`   //use provider to get default cluster slave dbs
			dbSlave1 Engine   `gone:"*,slave=0"` //use provider to get default cluster slave db 1
			dbSlave2 Engine   `gone:"*,slave=1"` //use provider to get default cluster slave db 2

			defaultCluster Engine   `gone:"xorm"`           // use provider to get default cluster db
			cluster2       Engine   `gone:"xorm,db=custom"` // use provider to get custom db
			cluster2Master Engine   `gone:"xorm,db=custom,master"`
			cluster2Slaves []Engine `gone:"xorm,db=custom,slave"`
			cluster2Slave1 Engine   `gone:"xorm,db=custom,slave=0"`
			cluster2Slave2 Engine   `gone:"xorm,db=custom,slave=1"`
		}) {
			assert.NotNil(t, in.db)
			assert.NotNil(t, in.dbMaster)
			assert.NotNil(t, in.dbSlaves)
			assert.NotNil(t, in.dbSlave1)
			assert.NotNil(t, in.dbSlave2)
			assert.NotNil(t, in.defaultCluster)
			assert.NotNil(t, in.cluster2)
			assert.NotNil(t, in.cluster2Master)
			assert.NotNil(t, in.cluster2Slaves)
			assert.NotNil(t, in.cluster2Slave1)
			assert.NotNil(t, in.cluster2Slave2)
		})
}
