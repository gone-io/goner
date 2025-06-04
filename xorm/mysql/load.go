package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/xorm"
)

// Load xorm.Engine and mysql driver
func Load(loader gone.Loader) error {
	loader.MustLoadX(xorm.Load)
	return nil
}
