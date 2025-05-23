package injector

import (
	"errors"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"math/rand"
	"reflect"
	"testing"
)

func Test_Prepare(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	type P struct {
		ID int
	}

	type C struct {
		p    P
		name string
	}

	executor := NewMockBindExecutor[*C](controller)

	executor.
		EXPECT().
		InjectedByType(gomock.Any()).
		DoAndReturn(func(pt reflect.Type) BindFunc[*C] {
			if reflect.TypeOf(P{}) != pt {
				return nil
			}

			return func(context *C) (reflect.Value, error) {
				return reflect.ValueOf(context.p), nil
			}
		}).
		Times(3)

	executor.
		EXPECT().
		FindFieldSetter(gomock.Any(), gomock.Any()).
		DoAndReturn(func(conf string, field reflect.StructField) (FieldSetter[*C], error) {
			return func(context *C, fieldValue reflect.Value) error {
				if conf == "" {
					if fieldValue.Type() == reflect.TypeOf(context.p) {
						fieldValue.Set(reflect.ValueOf(context.p))
					} else {
						fieldValue.Set(reflect.ValueOf(&context.p))
					}
				} else {
					fieldValue.Set(reflect.ValueOf(context.name))
				}
				return nil
			}, nil
		}).
		Times(3)

	gone.
		NewApp().
		Load(&delayBindInjector[*C]{name: "test-inject"}).
		Load(executor).
		Run(func(injector DelayBindInjector[*C], k gone.GonerKeeper, loader gone.Loader) {
			bindFn, err := injector.Prepare(func(p P, in struct {
				p1     P           `gone:"test-inject"`
				p2     *P          `gone:"test-inject"`
				name   string      `gone:"test-inject,name=food"`
				loader gone.Loader `gone:"*"`
			}, k gone.GonerKeeper) (p0, p1, p2 P, name string, keeper gone.GonerKeeper, loader gone.Loader) {
				return p, in.p1, *in.p2, in.name, k, in.loader
			})
			assert.Nil(t, err)

			var x = func() {
				c := C{
					p:    P{ID: rand.Int()},
					name: fmt.Sprintf("name-%d", rand.Int()),
				}

				result, err := bindFn(&c)
				assert.Nil(t, err)

				assert.Equal(t, P{ID: c.p.ID}, result[0].Interface())
				assert.Equal(t, P{ID: c.p.ID}, result[1].Interface())
				assert.Equal(t, P{ID: c.p.ID}, result[2].Interface())
				assert.Equal(t, c.name, result[3].Interface())
				assert.Equal(t, k, result[4].Interface())
				assert.Equal(t, loader, result[5].Interface())
			}

			for i := 0; i < 1000; i++ {
				go x()
			}
		})
}

func TestDelayBindInjector(t *testing.T) {
	type Ctx struct{}
	controller := gomock.NewController(t)
	defer controller.Finish()
	executor := NewMockBindExecutor[Ctx](controller)

	gone.NewApp().
		Load(executor).
		Load(&delayBindInjector[Ctx]{name: "x"}).
		Load(gone.WrapFunctionProvider(func(tagConf string, param struct{}) (*string, error) {
			var str = "hello"
			return &str, nil
		})).
		Run(func(injector DelayBindInjector[Ctx]) {
			t.Run("FindFieldSetter error", func(t *testing.T) {
				executor.EXPECT().FindFieldSetter(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
				executor.EXPECT().InjectedByType(gomock.Any()).Return(nil)
				prepare, err := injector.Prepare(func(i struct {
					x int `gone:"x,test"`
				}) {
				})
				assert.Nil(t, prepare)
				assert.Error(t, err)
			})
			t.Run("FindFieldSetter suc but parse err", func(t *testing.T) {
				executor.EXPECT().FindFieldSetter(gomock.Any(), gomock.Any()).Return(FieldSetter[Ctx](func(ctx Ctx, fieldValue reflect.Value) error {
					return errors.New("parse error")
				}), nil)
				executor.EXPECT().InjectedByType(gomock.Any()).Return(nil).AnyTimes()
				fn, err := injector.Prepare(func(i *struct {
					x int `gone:"x,test"`
				}) int {
					return i.x
				})
				assert.Nil(t, err)
				_, err = fn(Ctx{})
				assert.Error(t, err)
			})
			t.Run("success", func(t *testing.T) {
				executor.EXPECT().FindFieldSetter(gomock.Any(), gomock.Any()).Return(FieldSetter[Ctx](func(ctx Ctx, fieldValue reflect.Value) error {
					fieldValue.Set(reflect.ValueOf(1))
					return nil
				}), nil)
				executor.EXPECT().InjectedByType(gomock.Any()).Return(nil).AnyTimes()

				fn, err := injector.Prepare(func(i *struct {
					X int `gone:"x,test"`
				}, x *string) (int, string) {
					return i.X, *x
				})
				assert.Nil(t, err)
				values, err := fn(Ctx{})
				assert.Nil(t, err)
				assert.Equal(t, 1, values[0].Interface())
				assert.Equal(t, "hello", values[1].Interface())
			})

		})
}
