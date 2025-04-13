package xorm

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"xorm.io/xorm/log"
)

func Test_dbLogger_Level(t *testing.T) {
	logger := dbLogger{}

	logger.SetLevel(log.LOG_INFO)
	level := logger.Level()
	assert.Equal(t, log.LOG_INFO, level)

	logger.ShowSQL(false)
	assert.False(t, logger.IsShowSQL())
	logger.ShowSQL()
	assert.True(t, logger.IsShowSQL())
}

func Test_dbLogger(t *testing.T) {
	gone.
		Run(func(logger gone.Logger) {
			l := dbLogger{Logger: logger}
			l.Debug("print debug", "log")
			l.Info("print info", "log")
			l.Error("print error", "log")
			l.Warn("print warn", "log")
		})
}
