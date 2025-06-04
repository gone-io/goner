package apollo

import (
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/viper"
	"strings"
)

type changeListener struct {
	gone.Flag
	keyMap map[string][]any

	//lazy fill，resolve circular dependency
	logger gone.Logger `gone:"*" option:"lazy"`

	viper     *viper.Viper
	vipers    []*viper.Viper
	vipersMap map[string]*viper.Viper

	watchMap map[string]gone.ConfWatchFunc
}

func (c *changeListener) Put(key string, value any) {
	c.keyMap[key] = append(c.keyMap[key], value)
}

func (c *changeListener) Init() {
	c.vipersMap = make(map[string]*viper.Viper)
	c.vipers = make([]*viper.Viper, 0)
	c.keyMap = make(map[string][]any)
	c.watchMap = make(map[string]gone.ConfWatchFunc)
}

func (c *changeListener) SetViper(v *viper.Viper) {
	c.viper = v
}

func (c *changeListener) AddViper(namespace string, v *viper.Viper) {
	c.vipers = append(c.vipers, v)
	c.vipersMap[namespace] = v
}

func compare(a, b any) (equal bool) {
	return cmp.Equal(a, b)
}

func (c *changeListener) OnChange(event *storage.ChangeEvent) {
	oldValue := make(map[string]any)
	for key := range c.watchMap {
		oldValue[key] = c.viper.Get(key)
	}

	v := c.vipersMap[event.Namespace]
	if v != nil {
		for k, change := range event.Changes {
			switch change.ChangeType {
			case storage.DELETED:
				v.Set(k, nil)
			default:
				v.Set(k, change.NewValue)
			}
		}
	}

	c.viper = viper.New()
	for _, v := range c.vipers {
		err := c.viper.MergeConfigMap(v.AllSettings())
		g.ErrorPrinter(c.logger, err, "try to merge remote config")
	}

	for useKey, values := range c.keyMap {
		needUpdate := false
		for k := range event.Changes {
			if k == useKey || strings.HasPrefix(k, useKey) {
				needUpdate = true
				break
			}
		}

		if needUpdate {
			for _, ani := range values {
				err := c.viper.UnmarshalKey(useKey, ani)
				g.ErrorPrinter(c.logger, err, "try to unmarshal key")
			}
		}
	}
	for key, fn := range c.watchMap {
		oldVal := oldValue[key]
		newVal := c.viper.Get(key)
		if !compare(oldVal, newVal) {
			fn(oldVal, newVal)
		}
	}
}

func (c *changeListener) OnNewestChange(*storage.FullChangeEvent) {}

func (c *changeListener) Watch(key string, callback gone.ConfWatchFunc) {
	c.watchMap[key] = callback
}
