package sqlite

import (
	"github.com/gone-io/gone/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type dial struct {
	gone.Flag
	gorm.Dialector

	DriverName string `gone:"config,gorm.sqlite.driver-name"`
	DSN        string `gone:"config,gorm.sqlite.dsn"`
}

func (d *dial) Init() error {
	if d.Dialector == nil {
		d.Dialector = sqlite.New(sqlite.Config{
			DriverName: d.DriverName,
			DSN:        d.DSN,
		})
	}
	return nil
}

func Load(loader gone.Loader) error {
	return loader.Load(&dial{}, gone.IsDefault(new(gorm.Dialector)))
}
