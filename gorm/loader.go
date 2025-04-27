package gorm

import (
	"github.com/gone-io/gone/v2"
	"gorm.io/gorm/logger"
)

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&iLogger{}, gone.IsDefault(new(logger.Interface))).
		MustLoad(&dbProvider{})
	return nil
}
