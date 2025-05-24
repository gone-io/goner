package gin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin/injector"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type tester interface {
	TestF()
}
type xt struct {
}

func (t *xt) TestF() {
}

func Test_proxy(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	i := injector.NewMockDelayBindInjector[*gin.Context](controller)
	mockResponser := NewMockResponser(controller)

	p := &proxy{
		log:       gone.GetDefaultLogger(),
		responser: mockResponser,
		injector:  i,
		stat:      true,
	}

	t.Run("GonerName", func(t *testing.T) {
		assert.Equalf(t, IdGoneGinProxy, p.GonerName(), "GonerName()")
	})

	t.Run("Proxy", func(t *testing.T) {
		i.EXPECT().GonerName().Return("http")
		i.EXPECT().Prepare(gomock.Any()).Return(func(context *gin.Context) ([]reflect.Value, error) {
			return []reflect.Value{reflect.ValueOf("test=ok")}, nil
		}, nil).Times(2)

		gone.
			NewApp().
			Load(p).
			Load(i).
			Load(mockResponser).
			Run(func(p *proxy) {
				var handlers = []HandlerFunc{
					func(q Query[string]) string {
						return q.Get()
					},
					func(q Query[string]) string {
						return q.Get()
					},
				}
				arr := p.Proxy(handlers...)
				assert.Equal(t, 2, len(arr))
			})
	})

	t.Run("ProxyForMiddleware", func(t *testing.T) {
		i.EXPECT().GonerName().Return("http")
		i.EXPECT().Prepare(gomock.Any()).Return(func(context *gin.Context) ([]reflect.Value, error) {
			return []reflect.Value{reflect.ValueOf("test=ok")}, nil
		}, nil).Times(2)

		gone.
			NewApp().
			Load(p).
			Load(i).
			Load(mockResponser).
			Run(func(p *proxy) {
				var handlers = []HandlerFunc{
					func(q Query[string]) string {
						return q.Get()
					},
					func(q Query[string]) string {
						return q.Get()
					},
				}
				arr := p.ProxyForMiddleware(handlers...)
				assert.Equal(t, 2, len(arr))
			})
	})

	t.Run("proxyOne", func(t *testing.T) {
		type x struct {
		}

		f := func() *x {
			return nil
		}
		f2 := func() any {
			return (tester)(nil)
		}

		i.EXPECT().GonerName().Return("http")
		gone.
			NewApp().
			Load(p).
			Load(i).
			Load(mockResponser).
			Run(func(p *proxy) {
				p.stat = true

				t.Run("suc", func(t *testing.T) {
					x2 := &xt{}

					i.EXPECT().Prepare(gomock.Any()).Return(func(context *gin.Context) ([]reflect.Value, error) {
						return []reflect.Value{
							reflect.ValueOf(f2()),
							reflect.ValueOf(f()),
							reflect.ValueOf(x2),
							reflect.ValueOf("test=ok"),
						}, nil
					}, nil)
					var ctx = &gin.Context{}

					mockResponser.EXPECT().ProcessResults(ctx, nil, true, gomock.Any(), nil, nil, x2, "test=ok")
					fn := p.proxyOne(func(q Query[string]) (any, string) {
						return nil, q.Get()
					}, true)
					fn(ctx)
				})

				t.Run("err", func(t *testing.T) {
					err := errors.New("err")
					i.EXPECT().Prepare(gomock.Any()).Return(func(context *gin.Context) ([]reflect.Value, error) {
						return []reflect.Value{}, err
					}, nil)
					var ctx = &gin.Context{}

					mockResponser.EXPECT().Failed(ctx, err)
					fn := p.proxyOne(func(q Query[string]) (any, string) {
						return nil, q.Get()
					}, true)
					fn(ctx)
				})

			})
	})
}

func TestProxy(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockResponser := NewMockResponser(controller)

	type x struct {
	}

	var x1 = &x{}

	gone.
		NewApp(LoadGinHttpInjector).
		Load(&proxy{}).
		Load(mockResponser).
		Run(func(p *proxy) {

			fn := p.proxyOne(func(query Query[string], u *url.URL) (any, *x, string, *url.URL, error) {
				return x1, nil, query.Get(), u, nil
			}, true)

			addr, _ := url.Parse("http://localhost/?test=ok")

			ctx := &gin.Context{
				Request: &http.Request{
					URL: addr,
				},
			}
			mockResponser.EXPECT().ProcessResults(ctx, nil, true, gomock.Any(), x1, nil, "test=ok", addr, nil)
			fn(ctx)
		})
}
