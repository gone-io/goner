package gorm

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"gorm.io/gorm/logger"
)

var load = g.BuildOnceLoadFunc(
	g.L(&iLogger{}, gone.IsDefault(new(logger.Interface))),
	g.L(&dbProvider{}),
)

func Load(loader gone.Loader) error {
	return load(loader)
}
