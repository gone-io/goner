package goneMcp

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"strings"
)

type StdioConf struct {
	Command string   `json:"command"`
	Env     []string `json:"env"`
	Args    []string `json:"args"`
}

type SSEConf struct {
	BaseUrl string            `json:"baseUrl"`
	Header  map[string]string `json:"header"`
}

const paramKey = "param"
const configKey = "configKey"

func newStdioClient(m map[string]string, config gone.Configure) (c *client.Client, err error) {
	var conf StdioConf
	if m[configKey] != "" {
		if err = config.Get(m[configKey], &conf, ""); err != nil {
			return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("get mcp client config failed by key=%s", m[configKey]))
		}
	} else if m[paramKey] != "" {
		split := strings.Split(m[paramKey], " ")
		for _, it := range split {
			it = strings.TrimSpace(it)
			if it != "" {
				if conf.Command == "" {
					conf.Command = it
				} else {
					conf.Args = append(conf.Args, it)
				}
			}
		}
	}

	if conf.Command == "" {
		return nil, gone.ToError("command is empty")
	}

	c, err = client.NewStdioMCPClient(conf.Command, conf.Env, conf.Args...)
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("create mcp client failed by type=stdio"))
	}
	return
}

func newSseClient(m map[string]string, config gone.Configure) (c *client.Client, err error) {
	var conf SSEConf
	if m[configKey] != "" {
		if err = config.Get(m[configKey], &conf, ""); err != nil {
			return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("get mcp client config failed by key=%s", m[configKey]))
		}
	} else if m[paramKey] != "" {
		conf.BaseUrl = m[paramKey]
	}

	c, err = client.NewSSEMCPClient(conf.BaseUrl, client.WithHeaders(conf.Header))
	return c, gone.ToErrorWithMsg(err, fmt.Sprintf("create mcp client failed by type=sse"))
}

func newSseClientByInjectTransport(m map[string]string, keeper gone.GonerKeeper) (c *client.Client, err error) {
	trans, err := g.GetComponentByName[transport.Interface](keeper, m[paramKey])
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("can not found the transport by name=%s, please load it first", m[paramKey]))
	}
	c = client.NewClient(trans)
	return
}

// type=$[stdio|sse|transport],param=$[stdioCommand|sseBaseUrl,transportInjectedName],configKey=$[configKey]
var clientMap = make(map[string]*client.Client)

func clientProvide(
	tagConf string,
	param struct {
		keeper gone.GonerKeeper `gone:"*"`
		config gone.Configure   `gone:"configure"`
	},
) (c *client.Client, err error) {
	c = clientMap[tagConf]
	if c != nil {
		return c, nil
	}

	m, _ := gone.TagStringParse(tagConf)

	switch m["type"] {
	case "stdio":
		c, err = newStdioClient(m, param.config)
	case "sse":
		c, err = newSseClient(m, param.config)
	case "transport":
		c, err = newSseClientByInjectTransport(m, param.keeper)
	default:
		return nil, gone.ToError(fmt.Sprintf("support type=%s, inject config format: `gone:\"*,type=${stdio|sse|transport},param=${parameter},configKey=${configKey}\"`", m["type"]))
	}
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("create mcp client failed by type=%s", m["type"]))
	}
	clientMap[tagConf] = c
	return
}
