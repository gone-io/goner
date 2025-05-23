package injector

import (
	"errors"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

func Test_bindExecutor_Init(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	type Ctx struct {
		value string
	}

	parser := NewMockTypeParser[Ctx](controller)
	nameParser := NewMockNameParser[Ctx](controller)

	t.Run("duplicate type parser", func(t *testing.T) {
		parser.EXPECT().Type().Return(reflect.TypeOf("")).Times(2)
		err := gone.SafeExecute(func() error {
			gone.
				NewApp().
				Load(parser).
				Load(parser).
				Load(&bindExecutor[Ctx]{}).
				Run(func() {})
			return nil
		})
		assert.NotNil(t, err)
	})

	t.Run("duplicate name parser", func(t *testing.T) {
		nameParser.EXPECT().Name().Return("food").Times(2)
		err := gone.SafeExecute(func() error {
			gone.
				NewApp().
				Load(nameParser).
				Load(nameParser).
				Load(&bindExecutor[Ctx]{}).
				Run(func() {})
			return nil
		})
		assert.NotNil(t, err)
	})

	t.Run("InjectedByType", func(t *testing.T) {
		t.Run("InjectedByType not found", func(t *testing.T) {
			parser.EXPECT().Type().Return(reflect.TypeOf(""))
			gone.
				NewApp().
				Load(parser).
				Load(&bindExecutor[Ctx]{}).
				Run(func(executor BindExecutor[Ctx]) {
					fn := executor.InjectedByType(reflect.TypeOf(1))
					assert.Nil(t, fn)
				})
		})
		t.Run("InjectedByType success", func(t *testing.T) {
			parser.EXPECT().Type().Return(reflect.TypeOf(""))
			parser.EXPECT().Parse(gomock.Any()).Return(reflect.ValueOf("hello"), nil)
			gone.
				NewApp().
				Load(parser).
				Load(&bindExecutor[Ctx]{}).
				Run(func(executor BindExecutor[Ctx]) {
					fn := executor.InjectedByType(reflect.TypeOf(""))
					value, err := fn(Ctx{value: "hello"})
					assert.Nil(t, err)
					assert.Equal(t, "hello", value.Interface())
				})
		})
		t.Run("InjectedByType failed", func(t *testing.T) {
			parser.EXPECT().Type().Return(reflect.TypeOf(""))
			parser.EXPECT().Parse(gomock.Any()).Return(reflect.Value{}, errors.New("parse error"))
			gone.
				NewApp().
				Load(parser).
				Load(&bindExecutor[Ctx]{}).
				Run(func(executor BindExecutor[Ctx]) {
					fn := executor.InjectedByType(reflect.TypeOf(""))
					value, err := fn(Ctx{value: "hello"})
					assert.NotNil(t, err)
					assert.Equal(t, reflect.Value{}, value)
				})
		})
	})

	t.Run("FindFieldSetter", func(t *testing.T) {
		type Struct struct {
			V1 string `gone:"inject,conf"`
			V2 string `gone:"inject"`
		}
		field1, _ := reflect.TypeOf(Struct{}).FieldByName("V1")
		field2, _ := reflect.TypeOf(Struct{}).FieldByName("V2")
		//field3, _ := reflect.TypeOf(Struct{}).FieldByName("V3")

		t.Run("can not find parser", func(t *testing.T) {
			nameParser.EXPECT().Name().Return("confx")
			gone.
				NewApp().
				Load(nameParser).
				Load(&bindExecutor[Ctx]{}).
				Run(func(executor BindExecutor[Ctx]) {
					fn, err := executor.FindFieldSetter("conf", field1)
					assert.Nil(t, fn)
					assert.Error(t, err)
				})
		})
		t.Run("find name parser", func(t *testing.T) {
			nameParser.EXPECT().Name().Return("conf")
			gone.
				NewApp().
				Load(nameParser).
				Load(&bindExecutor[Ctx]{}).
				Run(func(executor BindExecutor[Ctx]) {

					t.Run("build parse err", func(t *testing.T) {
						nameParser.EXPECT().BuildParser(gomock.Any(), gomock.Any()).Return(nil, errors.New("build parse err"))
						fn, err := executor.FindFieldSetter("conf", field1)
						assert.Nil(t, fn)
						assert.Error(t, err)
					})
					t.Run("build parse ok", func(t *testing.T) {
						t.Run("parse err", func(t *testing.T) {
							nameParser.EXPECT().BuildParser(gomock.Any(), gomock.Any()).Return(func(ctx Ctx) (reflect.Value, error) {
								return reflect.Value{}, errors.New("parse err")
							}, nil)
							fn, err := executor.FindFieldSetter("conf", field1)
							assert.Nil(t, err)
							var value reflect.Value
							err = fn(Ctx{}, value)
							assert.Error(t, err)
						})
						t.Run("parse ok", func(t *testing.T) {
							nameParser.EXPECT().BuildParser(gomock.Any(), gomock.Any()).Return(func(ctx Ctx) (reflect.Value, error) {
								return reflect.ValueOf("hello"), nil
							}, nil)
							fn, err := executor.FindFieldSetter("conf", field1)
							assert.Nil(t, err)
							var value = reflect.New(reflect.TypeOf(""))
							err = fn(Ctx{}, value.Elem())
							assert.Nil(t, err)
							assert.Equal(t, "hello", value.Elem().Interface())
						})

					})
				})
		})

		t.Run("can not find parser for type", func(t *testing.T) {
			gone.
				NewApp().
				Load(&bindExecutor[Ctx]{}).
				Run(func(executor BindExecutor[Ctx]) {
					fn, err := executor.FindFieldSetter("", field2)
					assert.Nil(t, fn)
					assert.Error(t, err)
				})
		})

		t.Run("find type parser", func(t *testing.T) {
			parser.EXPECT().Type().Return(reflect.TypeOf(""))
			gone.
				NewApp().
				Load(parser).
				Load(&bindExecutor[Ctx]{}).
				Run(func(executor BindExecutor[Ctx]) {
					t.Run("parse success", func(t *testing.T) {
						parser.EXPECT().Parse(gomock.Any()).Return(reflect.ValueOf("hello"), nil)
						fn, err := executor.FindFieldSetter("", field2)
						assert.Nil(t, err)
						var value = reflect.New(reflect.TypeOf(""))
						err = fn(Ctx{}, value.Elem())
						assert.Nil(t, err)
						assert.Equal(t, "hello", value.Elem().Interface())
					})
					t.Run("parse failed", func(t *testing.T) {
						parser.EXPECT().Parse(gomock.Any()).Return(reflect.Value{}, errors.New("parse err"))
						fn, err := executor.FindFieldSetter("", field2)
						assert.Nil(t, err)
						var value = reflect.New(reflect.TypeOf(""))
						err = fn(Ctx{}, value.Elem())
						assert.Error(t, err)
					})

				})
		})
	})
}
