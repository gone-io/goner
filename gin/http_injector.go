package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin/injector"
	"github.com/gone-io/goner/gin/parser"
)

// LoadGinHttpInjector load http injector
func LoadGinHttpInjector(loader gone.Loader) error {
	loader.
		MustLoadX(injector.BuildLoad[*gin.Context](IdHttpInjector)).
		MustLoadX(parser.Load)
	return nil
}
