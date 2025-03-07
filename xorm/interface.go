package xorm

import (
	"xorm.io/xorm"
)

type Interface = xorm.Interface

type XormEngine interface {
	xorm.EngineInterface
	Transaction(fn func(session xorm.Interface) error) error
	Sqlx(sql string, args ...any) *xorm.Session
	GetOriginEngine() xorm.EngineInterface
	SetPolicy(policy xorm.GroupPolicy)
}

type Engine = XormEngine
