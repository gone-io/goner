package gin_test

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

//go:generate mockgen -package=gin_test -destination ./tracer_mock_test.go github.com/gone-io/goner/g Tracer

func Test_proxy_Proxy(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	responser := gin.NewMockResponser(controller)
	injector := gin.NewMockHttInjector(controller)
	responser.EXPECT().Success(gomock.Any(), gomock.Any()).AnyTimes()
	responser.EXPECT().ProcessResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.All()).AnyTimes()

	//mockTracer := NewMockTracer(controller)
	//mockTracer.EXPECT().GetTraceId().Return("123")

	gone.
		NewApp(func(cemetery gone.Loader) error {
			if err := cemetery.Load(gin.NewGinProxy()); err != nil {
				return err
			}
			if err := cemetery.Load(responser); err != nil {
				return err
			}
			if err := cemetery.Load(injector); err != nil {
				return err
			}
			//if err := cemetery.Load(mockTracer); err != nil {
			//	return err
			//}
			return nil
		}).
		Test(func(proxy gin.HandleProxyToGin, logger gone.Logger) {
			i := 0
			t.Run("Special funcs", func(t *testing.T) {
				funcs := proxy.Proxy(
					func(*gin.Context) (any, error) {
						i++
						return nil, nil
					},
					func(*gin.Context) error {
						i++
						return nil
					},
					func(*gin.Context) {
						i++
					},
					func(*gin.OriginContent) (any, error) {
						i++
						return nil, nil
					},
					func(*gin.OriginContent) error {
						i++
						return nil
					},
					func(*gin.OriginContent) {
						i++
					},
					func() {
						i++
					},
					func() (any, error) {
						i++
						return nil, nil
					},
					func() error {
						i++
						return nil
					},
				)
				for _, fn := range funcs {
					fn(&gin.OriginContent{})
				}

				assert.Equal(t, 9, i)
			})

			t.Run("Inject Error", func(t *testing.T) {
				defer func() {
					err := recover()
					assert.NotNil(t, err)
				}()

				injector.EXPECT().StartBindFuncs()

				proxy.ProxyForMiddleware(func(in struct {
					x gone.Logger `gone:"xxx"`
				}) {

				})

			})

			t.Run("Bind Context Error", func(t *testing.T) {
				bindErr := errors.New("bind error")

				injector.EXPECT().StartBindFuncs()
				injector.EXPECT().BindFuncs().Return(func(ctx *gin.OriginContent, obj reflect.Value) (reflect.Value, error) {
					return reflect.Value{}, bindErr
				})

				responser.EXPECT().Failed(gomock.Any(), gomock.Any()).Do(func(ctx any, err error) {
					assert.Equal(t, bindErr, err)
				})

				arr := proxy.ProxyForMiddleware(func(in struct {
					x gone.Logger `gone:"*"`
				}) {
				})
				arr[0](&gin.OriginContent{})
			})
		})
}
