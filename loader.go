package goner

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin"
	"github.com/gone-io/goner/tracer"
	viper "github.com/gone-io/goner/viper"
	zap "github.com/gone-io/goner/zap"
)

func BaseLoad(loader gone.Loader) error {
	return gone.OnceLoad(func(loader gone.Loader) error {
		err := tracer.Load(loader)
		if err != nil {
			return err
		}
		err = viper.Load(loader)
		if err != nil {
			return err
		}
		return zap.Load(loader)
	})(loader)
}

func GinLoad(loader gone.Loader) error {
	return gone.OnceLoad(func(loader gone.Loader) error {
		if err := BaseLoad(loader); err != nil {
			return gone.ToError(err)
		}
		return gin.Load(loader)
	})(loader)
}
