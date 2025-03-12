package urllib

import (
	"github.com/gone-io/gone/v2"
	"github.com/imroc/req/v3"
)

type requestProvider struct {
	gone.Flag
	client Client `gone:"*"`
}

func (p *requestProvider) Provide() (*req.Request, error) {
	return p.client.R(), nil
}

type clientProvider struct {
	gone.Flag
	client Client `gone:"*"`
}

func (p *clientProvider) Provide() (*req.Client, error) {
	return p.client.C(), nil
}
