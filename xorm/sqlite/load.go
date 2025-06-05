package sqlite

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/xorm"
	_ "github.com/mattn/go-sqlite3"
)

// Load xorm.Engine and sqlite3 driver
func Load(loader gone.Loader) error {
	loader.MustLoadX(xorm.Load)
	return nil
}
