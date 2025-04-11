package xorm

import (
	"errors"
	mock "github.com/gone-io/gone/mock/v2"
	"go.uber.org/mock/gomock"
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
