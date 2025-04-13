package xorm

import (
	"github.com/gone-io/gone/v2"
	"xorm.io/xorm"
)

func newEng(xEng xorm.EngineInterface, logger gone.Logger) *eng {
	e := eng{EngineInterface: xEng}
	e.trans = newTrans(logger, func() Session {
		return e.NewSession()
	})
	return &e
}

type eng struct {
	xorm.EngineInterface
	trans
}

func (e *eng) GetOriginEngine() xorm.EngineInterface {
	return e.EngineInterface
}

func (e *eng) SetPolicy(policy xorm.GroupPolicy) {
	if group, ok := e.EngineInterface.(*xorm.EngineGroup); ok {
		group.SetPolicy(policy)
	}
}

func (e *eng) Sqlx(sql string, args ...any) *xorm.Session {
	sql, args = sqlDeal(sql, args...)
	return e.SQL(sql, args...)
}
