package viper

import (
	"github.com/go-viper/encoding/javaproperties"
	"github.com/gone-io/gone/v2"
	"github.com/spf13/afero"

	"github.com/spf13/viper"
	"reflect"
	"strings"
)

func New(testFlag gone.TestFlag) gone.Configure {
	return &configure{testFlag: testFlag}
}

type configure struct {
	gone.Flag
	testFlag gone.TestFlag `gone:"*" option:"allowNil"`
	conf     *viper.Viper
}

func (c *configure) readConfig() (err error) {
	codecRegistry := viper.NewCodecRegistry()
	codec := &javaproperties.Codec{}
	err = codecRegistry.RegisterCodec("properties", codec)
	if err != nil {
		return gone.ToError(err)
	}

	conf := viper.NewWithOptions(viper.WithCodecRegistry(codecRegistry))
	conf.SetEnvPrefix("gone")
	conf.AutomaticEnv()
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if files, err := getConfigFiles(c.testFlag != nil, afero.NewOsFs()); err != nil {
		return gone.ToError(err)
	} else {
		for _, f := range files {
			conf.SetConfigFile(f)
			conf.SetConfigType(fileExt(f))
			err = conf.MergeInConfig()
			if err != nil {
				return gone.ToError(err)
			}
		}
	}
	c.conf = conf
	return
}

func (c *configure) Get(key string, v any, defaultVal string) error {
	if c.conf == nil {
		err := c.readConfig()
		if err != nil {
			return err
		}
	}
	return c.get(key, v, defaultVal)
}

func (c *configure) get(key string, value any, defaultVale string) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return gone.NewInnerError("type of value must be ptr", gone.NotSupport)
	}
	v := c.conf.Get(key)
	if v == nil || v == "" {
		return gone.SetValue(rv, value, defaultVale)
	}
	return gone.ToError(c.conf.UnmarshalKey(key, value))
}
