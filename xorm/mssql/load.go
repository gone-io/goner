package mssql

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/xorm"
)

// Load xorm.Engine and mssql driver
func Load(loader gone.Loader) error {
	loader.MustLoadX(xorm.Load)
	return nil
}
