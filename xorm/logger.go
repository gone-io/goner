package xorm

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"xorm.io/xorm/log"
)

type dbLogger struct {
	gone.Logger
	showSql bool
	level   log.LogLevel
}

func (l *dbLogger) Level() log.LogLevel {
	return l.level
}
func (l *dbLogger) SetLevel(level log.LogLevel) {
	l.level = level
}

func (l *dbLogger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		l.showSql = show[0]
	} else {
		l.showSql = true
	}
}
func (l *dbLogger) IsShowSQL() bool {
	return l.showSql
}

func (l *dbLogger) Debug(v ...interface{}) {
	l.Logger.Debugf("%s", fmt.Sprintln(v...))
}
func (l *dbLogger) Error(v ...interface{}) {
	l.Logger.Errorf("%s", fmt.Sprintln(v...))
}
func (l *dbLogger) Info(v ...interface{}) {
	l.Logger.Infof("%s", fmt.Sprintln(v...))
}
func (l *dbLogger) Warn(v ...interface{}) {
	l.Logger.Warnf("%s", fmt.Sprintln(v...))
}
