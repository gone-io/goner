package goneMcp

import (
	"fmt"
	"time"

	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/mark3labs/mcp-go/server"
)

type SSEConfig struct {
	WithBaseURL                      string        `json:"withBaseURL,omitempty"`
	WithBasePath                     string        `json:"withBasePath,omitempty"`
	WithMessageEndpoint              string        `json:"withMessageEndpoint,omitempty"`
	WithUseFullURLForMessageEndpoint bool          `json:"withUseFullURLForMessageEndpoint,omitempty"`
	WithSSEEndpoint                  string        `json:"withSSEEndpoint,omitempty"`
	WithKeepAliveInterval            time.Duration `json:"withKeepAliveInterval,omitempty"`
	WithKeepAlive                    bool          `json:"withKeepAlive,omitempty"`
	Address                          string        `json:"address"`
}

func (conf SSEConfig) ToOptions() []server.SSEOption {
	var options []server.SSEOption
	if conf.WithBaseURL != "" {
		options = append(options, server.WithBaseURL(conf.WithBaseURL))
	}
	if conf.WithBasePath != "" {
		options = append(options, server.WithBasePath(conf.WithBasePath))
	}
	if conf.WithMessageEndpoint != "" {
		options = append(options, server.WithMessageEndpoint(conf.WithMessageEndpoint))
	}
	if conf.WithUseFullURLForMessageEndpoint {
		options = append(options, server.WithUseFullURLForMessageEndpoint(true))
	}
	if conf.WithSSEEndpoint != "" {
		options = append(options, server.WithSSEEndpoint(conf.WithSSEEndpoint))
	}
	if conf.WithKeepAliveInterval > 0 {
		options = append(options, server.WithKeepAliveInterval(conf.WithKeepAliveInterval))
	}
	if conf.WithKeepAlive {
		options = append(options, server.WithKeepAlive(true))
	}
	return options
}

type Conf struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`

	WithRecovery           bool   `json:"withRecovery,omitempty"`
	WithPromptCapabilities bool   `json:"withPromptCapabilities,omitempty"`
	WithToolCapabilities   bool   `json:"withToolCapabilities,omitempty"`
	WithLogging            bool   `json:"withLogging,omitempty"`
	WithInstructions       string `json:"withInstructions,omitempty"`
	TransportType          string `json:"transportType"`
}

func (conf Conf) ToOptions() []server.ServerOption {
	var options []server.ServerOption
	if conf.WithRecovery {
		options = append(options, server.WithRecovery())
	}
	if conf.WithPromptCapabilities {
		options = append(options, server.WithPromptCapabilities(true))
	}
	if conf.WithToolCapabilities {
		options = append(options, server.WithToolCapabilities(true))
	}
	if conf.WithLogging {
		options = append(options, server.WithLogging())
	}
	if conf.WithInstructions != "" {
		options = append(options, server.WithInstructions(conf.WithInstructions))
	}
	return options
}

type serverProvider struct {
	gone.Flag

	config      gone.Configure   `gone:"configure"`
	logger      gone.Logger      `gone:"*"`
	keeper      gone.GonerKeeper `gone:"*"`
	beforeStart gone.BeforeStart `gone:"*"`
	beforeStop  gone.BeforeStop  `gone:"*"`

	tools     []ITool     `gone:"*"`
	props     []IPrompt   `gone:"*"`
	resources []IResource `gone:"*"`

	m map[string]*server.MCPServer
}

func (s *serverProvider) Init() error {
	s.m = make(map[string]*server.MCPServer)

	if len(s.tools) > 0 || len(s.props) > 0 || len(s.resources) > 0 {
		if mcpServer, err := s.Provide(""); err != nil {
			return gone.ToError(err)
		} else {
			for _, tool := range s.tools {
				mcpServer.AddTool(tool.Define(), tool.Process())
			}
			for _, prompt := range s.props {
				mcpServer.AddPrompt(prompt.Define(), prompt.Process())
			}
			for _, resource := range s.resources {
				mcpServer.AddResource(resource.Define(), resource.Process())
			}
		}
	}
	return nil
}

func (s *serverProvider) Provide(tagConf string) (*server.MCPServer, error) {
	_, keys := gone.TagStringParse(tagConf)
	var key = "mcp"
	if len(keys) > 0 && keys[0] != "" {
		key = keys[0]
	}
	mcpServer := s.m[key]
	if mcpServer != nil {
		return mcpServer, nil
	}

	var conf Conf
	if err := s.config.Get(key, &conf, ""); err != nil {
		return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("get mcp server config failed by key=%s", key))
	}

	options := conf.ToOptions()
	if hooks, err := g.GetComponentByName[*server.Hooks](s.keeper, key+".hooks"); err == nil {
		options = append(options, server.WithHooks(hooks))
	}

	mcpServer = server.NewMCPServer(conf.Name, conf.Version, options...)

	switch conf.TransportType {
	case "sse":
		var sseConf SSEConfig
		sseKey := key + ".sse"
		if err := s.config.Get(sseKey, &sseConf, ""); err != nil {
			return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("get mcp sse server config failed by key=%s", sseKey))
		}
		ops := sseConf.ToOptions()
		if fn, err := g.GetComponentByName[server.SSEContextFunc](s.keeper, sseKey+".context"); err == nil {
			ops = append(ops, server.WithSSEContextFunc(fn))
		}

		sseServer := server.NewSSEServer(mcpServer, ops...)

		s.beforeStart(func() {
			address := sseConf.Address
			if sseConf.Address == "" {
				address = ":8080"
			}
			s.logger.Infof("mcp: start sse server at %s", address)
			go func() {
				if err := sseServer.Start(address); err != nil {
					s.logger.Errorf("mcp: start sse server err: %v", err)
				}
			}()
		})

		// there is a data trace for sseServer.Start writing `srv` and sseServer.Shutdown read `srv`,
		// look for https://github.com/mark3labs/mcp-go/issues/166
		//s.beforeStop(func() {
		//	if err := sseServer.Shutdown(context.Background()); err != nil {
		//		s.logger.Errorf("mcp: shutdown sse server err: %v", err)
		//	}
		//})

		s.m[key] = mcpServer
		return mcpServer, nil

	case "", "stdio":
		var ops []server.StdioOption
		if fn, err := g.GetComponentByName[server.StdioContextFunc](s.keeper, key+".stdio.context"); err == nil {
			ops = append(ops, server.WithStdioContextFunc(fn))
		}

		s.beforeStart(func() {
			go func() {
				s.logger.Infof("mcp: start stdio server")
				if err := server.ServeStdio(mcpServer, ops...); err != nil {
					s.logger.Errorf("mcp: serve stdio err: %v", err)
				}
			}()
		})

		s.m[key] = mcpServer
		return mcpServer, nil
	default:
		return nil, gone.ToError("unsupported TransportType, only support sse and stdio.")
	}
}
