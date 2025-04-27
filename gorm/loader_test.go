package gorm

import (
	"github.com/gone-io/goner/g"
	"gorm.io/gorm"
	"testing"

	"github.com/gone-io/gone/v2"
	"go.uber.org/mock/gomock"
)

func TestLoad_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dialector := NewMockDialector(ctrl)

	gone.NewApp(Load, g.NamedThirdComponentLoadFunc[gorm.Dialector]("", dialector)).Run()
}
