package clickhouse

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestPriest(t *testing.T) {
	gone.RunTest(func(in struct {
		dial gorm.Dialector `gone:"*"`
	}) {
		assert.NotNil(t, in.dial)
	}, Load)
}
