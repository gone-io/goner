package xorm

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_eng_Sqlx(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	engineInterface := NewMockEngineInterface(controller)

	logger := gone.GetDefaultLogger()

	e := newEng(engineInterface, logger)
	assert.Equal(t, engineInterface, e.GetOriginEngine())

	engineInterface.EXPECT().SQL("select * from user where id = ?", 1).Return(nil)

	sqlx := e.Sqlx("select * from user where id = ?", 1)
	assert.Nil(t, sqlx)
}
