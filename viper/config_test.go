package viper

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfigure_Get(t *testing.T) {
	v := viper.New()
	v.Set("test.empty1", "")
	v.Set("test.int", 100)
	v.Set("test.string", "hello")
	v.Set("test.bool", true)
	v.Set("test.float", 1.2)
	v.Set("test.duration", "1s")
	v.Set("test.time", "2022-01-01 00:00:00")
	v.Set("test.intSlice", []int{1, 2, 3})
	v.Set("test.stringSlice", []string{"a", "b", "c"})
	v.Set("test.map", map[string]any{
		"a": 1,
		"b": "2",
		"c": true,
	})

	type Test struct {
		Int         int
		String      string
		Bool        bool
		Float       float64
		Duration    time.Duration
		Time        time.Time
		IntSlice    []int
		StringSlice []string
		Map         map[string]any
	}

	t.Run("use default value when empty", func(t *testing.T) {
		c := configure{
			conf: v,
		}
		var test int

		err := c.get("test.empty", &test, "900")
		assert.Nil(t, err)
		assert.Equal(t, 900, test)

		var test1 int
		err = c.get("test.empty1", &test1, "900")
		assert.Nil(t, err)
		assert.Equal(t, 900, test1)
	})

	t.Run("v is struct", func(t *testing.T) {
		var test struct {
			A int
			B string
			C bool
		}
		c := configure{
			conf: v,
		}
		err := c.get("test.map", &test, "900")
		assert.Nil(t, err)
		assert.Equal(t, 1, test.A)
		assert.Equal(t, "2", test.B)
		assert.Equal(t, true, test.C)
	})

	t.Run("v is custom type", func(t *testing.T) {
		type CustomType int
		var test CustomType
		c := configure{
			conf: v,
		}
		err := c.get("test.int", &test, "900")
		assert.Nil(t, err)
		assert.Equal(t, CustomType(100), test)
	})
}

func Test_configure_readConfig(t *testing.T) {
	c := configure{}
	err := c.readConfig()
	assert.Nil(t, err)
}

func TestNew(t *testing.T) {
	g := New(nil)
	assert.NotNil(t, g)
}
