package dependent

import "github.com/gone-io/gone/v2"

type iDependent struct {
	gone.Flag
}

func (s *iDependent) DoSomething() error {
	return nil
}
