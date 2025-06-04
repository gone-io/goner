package postgres

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/xorm"
	_ "github.com/lib/pq"
)

// Load xorm.Engine and postgres driver.
func Load(loader gone.Loader) error {
	loader.MustLoadX(xorm.Load)
	return nil
}
