package apollo

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/gone-io/gone/v2"
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
}

func (c *changeListener) Put(key string, value any) {
	c.keyMap[key] = append(c.keyMap[key], value)
}

func (c *changeListener) Init() {
	c.vipersMap = make(map[string]*viper.Viper)
	c.vipers = make([]*viper.Viper, 0)
	c.keyMap = make(map[string][]any)
}

func (c *changeListener) SetViper(v *viper.Viper) {
	c.viper = v
}

func (c *changeListener) AddViper(namespace string, v *viper.Viper) {
	c.vipers = append(c.vipers, v)
	c.vipersMap[namespace] = v
}

func (c *changeListener) OnChange(event *storage.ChangeEvent) {
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

	for _, v := range c.vipers {
		err := c.viper.MergeConfigMap(v.AllSettings())
		if err != nil {
			c.logger.Warnf("try to merge remote config err:%v\n", err)
		}
	}

	x := c.viper.Get("test.key")
	fmt.Printf("%v", x)

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
				if err := c.viper.UnmarshalKey(useKey, ani); err != nil {
					c.logger.Warnf("try to set `%s` value err:%v\n", useKey, err)
				}
			}
		}
	}
}

func (c *changeListener) OnNewestChange(*storage.FullChangeEvent) {}
