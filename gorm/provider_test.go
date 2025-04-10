package gorm

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gone-io/gone/v2"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

//go:generate mockgen -package gorm -destination=./dialector_mock_test.go gorm.io/gorm Dialector,ConnPool
//go:generate mockgen -package gorm -destination=./gorm_logger_mock_test.go gorm.io/gorm/logger Interface

func Test_dbProvider_Provide(t *testing.T) {
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDb.Close()

	mock.ExpectBegin()

	controller := gomock.NewController(t)
	defer controller.Finish()

	dialector := NewMockDialector(controller)
	dialector.EXPECT().Initialize(gomock.Any()).Return(nil)

	provider := gone.WrapFunctionProvider(func(string2 string, in struct{}) (gorm.Dialector, error) {
		return dialector, nil
	})

	gone.
		NewApp(Load).
		Load(provider).
		Run(func(p *dbProvider) {
			_, err := p.Provide("")
			if err == nil {
				t.Fatal(err)
			}

			dialector.EXPECT().Initialize(gomock.Any()).DoAndReturn(func(db *gorm.DB) error {
				db.ConnPool = sqlDb
				return nil
			})

			p.MaxIdle = 1
			p.MaxOpen = 1
			_, err = p.Provide("")
			if err != nil {
				t.Fatal(err)
			}

			_, err = p.Provide("")
			if err != nil {
				t.Fatal(err)
			}
		})
}
