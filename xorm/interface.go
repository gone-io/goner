package xorm

import (
	"database/sql/driver"
	"github.com/gone-io/gone/v2"
	"io"
	"xorm.io/xorm"
)

type Interface = xorm.Interface

// Engine db engine
type Engine interface {
	xorm.EngineInterface
	Transaction(fn func(session xorm.Interface) error) error
	Sqlx(sql string, args ...any) *xorm.Session
	GetOriginEngine() xorm.EngineInterface
	SetPolicy(policy xorm.GroupPolicy)
}

type Session interface {
	xorm.Interface
	driver.Tx
	io.Closer
	Begin() error
}
type XInterface = Session

// XormEngine @Deprecated use Engine
type XormEngine = Engine

var xormInterface = gone.GetInterfaceType(new(XormEngine))
var xormInterfaceSlice = gone.GetInterfaceType(new([]XormEngine))
