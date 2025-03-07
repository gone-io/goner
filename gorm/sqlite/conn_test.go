package sqlite

import (
	"github.com/gone-io/gone/v2"
	gone_viper "github.com/gone-io/goner/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestPriest(t *testing.T) {
	gone.RunTest(func(in struct {
		dial gorm.Dialector `gone:"*"`
	}) {
		assert.NotNil(t, in.dial)
	}, Load, gone_viper.Load)
}
