package injector

import (
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin/injector/parser"
)

func LoadGinHttpInjector(loader gone.Loader) error {
	loader.
		MustLoad(&delayBindInjector[*gin.Context]{name: "http"}).
		MustLoad(&bindExecutor[*gin.Context]{}).
		MustLoadX(parser.Load)
	return nil
}
