package gorm

import (
	"database/sql"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

//go:generate mockgen -package gorm -destination=./dialector_mock_test.go gorm.io/gorm Dialector
//go:generate mockgen -package gorm -destination=./pool_mock_test.go -source=./priest_test.go

type TestPool interface {
	gorm.GetDBConnector
	gorm.ConnPool
}

func TestPriest(t *testing.T) {
	controller := gomock.NewController(t)
	dialector := NewMockDialector(controller)
	pool := NewMockTestPool(controller)
	db := sql.DB{}
	pool.EXPECT().GetDBConn().Return(&db, nil)

	dialector.EXPECT().Initialize(gomock.Any()).DoAndReturn(func(db *gorm.DB) error {
		db.ConnPool = pool
		return nil
	})

	gone.RunTest(func(
		in struct {
			db *gorm.DB `gone:"*"`
		},
	) {
		assert.NotNil(t, in.db)
	}, func(loader gone.Loader) error {
		err := loader.Load(gone.WrapFunctionProvider(func(tagConf string, p struct{}) (gorm.Dialector, error) { return dialector, nil }))
		assert.Nil(t, err)
		if err != nil {
			return err
		}
		return Priest(loader)
	})
}
