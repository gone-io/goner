package module

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/xorm"
	"template_module/internal/interface/entity"
	"template_module/internal/module/user"
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
