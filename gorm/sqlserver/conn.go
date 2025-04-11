package sqlserver

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type dial struct {
	gone.Flag
	gorm.Dialector

	DriverName        string `gone:"config,gorm.sqlserver.driver-name"`
	DSN               string `gone:"config,gorm.sqlserver.dsn"`
	DefaultStringSize int    `gone:"config,gorm.sqlserver.default-string-size"`
}

func (d *dial) Init() error {
	if d.Dialector == nil {
		d.Dialector = sqlserver.New(sqlserver.Config{
			DriverName:        d.DriverName,
			DSN:               d.DSN,
			DefaultStringSize: d.DefaultStringSize,
		})
	}
	return nil
}

var load = g.BuildOnceLoadFunc(g.L(&dial{}, gone.IsDefault(new(gorm.Dialector))))

func Load(loader gone.Loader) error {
	return load(loader)
}
