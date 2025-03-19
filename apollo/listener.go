package apollo

import (
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/gone-io/gone/v2"
)

type changeListener struct {
	gone.Flag
	keyMap map[string][]any

	//lazy fill，resolve circular dependency
	logger gone.Logger `gone:"*" option:"lazy"`
}

func (c *changeListener) Init() {
	c.keyMap = make(map[string][]any)
}

func (c *changeListener) add(key string, v any) {
	c.keyMap[key] = append(c.keyMap[key], v)
}

func (c *changeListener) OnChange(event *storage.ChangeEvent) {
	for k, change := range event.Changes {
		if w, ok := c.keyMap[k]; ok && change.ChangeType == storage.MODIFIED {
			for _, v := range w {
				err := setValue(v, change.NewValue)
				if err != nil && c.logger != nil {
					c.logger.Warnf("try to change `%s` value  err: %v\n", k, err)
				}
			}
		}
	}
}

func (c *changeListener) OnNewestChange(*storage.FullChangeEvent) {}
