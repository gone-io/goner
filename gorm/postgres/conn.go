package postgres

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dial struct {
	gone.Flag
	gorm.Dialector

	driverName           string `gone:"config,gorm.postgres.driver-name"`
	dsn                  string `gone:"config,gorm.postgres.dsn"`
	withoutQuotingCheck  bool   `gone:"config,gorm.postgres.without-quoting-check,default=false"`
	preferSimpleProtocol bool   `gone:"config,gorm.postgres.prefer-simple-protocol,default=false"`
	withoutReturning     bool   `gone:"config,gorm.postgres.without-returning=false"`
}

func (d *dial) Init() error {
	if d.Dialector == nil {
		d.Dialector = postgres.New(postgres.Config{
			DriverName:           d.driverName,
			DSN:                  d.dsn,
			WithoutReturning:     d.withoutReturning,
			PreferSimpleProtocol: d.preferSimpleProtocol,
			WithoutQuotingCheck:  d.withoutQuotingCheck,
		})
	}
	return nil
}

var load = g.BuildOnceLoadFunc(g.L(&dial{}, gone.IsDefault(new(gorm.Dialector))))

func Load(loader gone.Loader) error {
	return load(loader)
}
