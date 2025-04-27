package user

import (
	"examples/gin_xorm_viper/internal"
	"examples/gin_xorm_viper/internal/interface/entity"
	"examples/gin_xorm_viper/internal/interface/mock"
	_ "github.com/go-sql-driver/mysql" //导入mysql驱动
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_iUser_Register(t *testing.T) {
	gone.RunTest(func(in struct {
		iUser      *iUser               `gone:"*"` //inject iUser for test
		iDependent *mock.MockIDependent `gone:"*"` //inject iDependent for mock
	}) {
		err := gone.ToError("err")
		in.iDependent.EXPECT().DoSomething().Return(err)

		register, err2 := in.iUser.Register(&entity.RegisterParam{
			Username: "test",
			Password: "test",
		})
		assert.Nil(t, register)
		assert.Equal(t, err2, err)
	}, func(loader gone.Loader) error {
		controller := gomock.NewController(t)

		//load all mocked components
		mock.MockLoader(loader, controller)

		err := loader.Load(&iUser{})
		if err != nil {
			return gone.ToError(err)
		}

		return internal.TestLoader(loader, controller)
	})
}
