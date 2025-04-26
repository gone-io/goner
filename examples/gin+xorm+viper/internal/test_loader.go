package internal

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/xorm"
	"go.uber.org/mock/gomock"
)

func TestLoader(loader gone.Loader, ctrl *gomock.Controller) error {
	return xorm.Load(loader)
}
