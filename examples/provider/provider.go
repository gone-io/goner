package provider

import "github.com/gone-io/gone/v2"

var _ gone.Provider[*ThirdComponent1] = (*provider)(nil)

type provider struct {
	gone.Flag
}

func (p *provider) Provide(tagConf string) (*ThirdComponent1, error) {
	//TODO implement me
	panic("implement me")
}

var _ gone.NoneParamProvider[*ThirdComponent2] = (*noneParamProvider)(nil)

type noneParamProvider struct {
	gone.Flag
}

func (p noneParamProvider) Provide() (*ThirdComponent2, error) {
	//TODO implement me
	panic("implement me")
}
