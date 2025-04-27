package module

import (
	"examples/gin_xorm_viper/internal/interface/entity"
	"examples/gin_xorm_viper/internal/module/user"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/xorm"
)

type dbSys struct {
	gone.Flag
	db xorm.Engine `gone:"*"`
}

func (s *dbSys) Init() error {
	return s.db.Sync(
		new(entity.User),
		new(user.TokenRecord),
	)
}
