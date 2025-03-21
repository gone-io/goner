package viper

import (
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	gone.
		NewApp(Load).
		Test(func(in struct {
			test int `gone:"config,test.ini"`
		}) {
			assert.Equal(t, 9900, in.test)
		})
}
