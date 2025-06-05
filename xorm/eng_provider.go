package xorm

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/spf13/cast"
	"reflect"
)

type engProvider struct {
	gone.Flag
	logger    gone.Logger   `gone:"*"`
	xProvider *xormProvider `gone:"*"`
}

func (s *engProvider) GonerName() string {
	return "xorm"
}

func (s *engProvider) Provide(tagConf string, t reflect.Type) (any, error) {
	switch t {
	case xormInterfaceSlice:
		group, err := s.xProvider.ProvideEngineGroup(tagConf)
		if err != nil {
			return nil, gone.ToError(err)
		}
		slaves := group.Slaves()
		engines := make([]Engine, 0, len(slaves))
		for _, slave := range slaves {
			engines = append(engines, newEng(slave, s.logger))
		}
		return engines, nil
	case xormInterface:
		return s.ProvideEngine(tagConf)

	default:
		return nil, gone.NewInnerErrorWithParams(gone.GonerTypeNotMatch, "Cannot find matched value for %q", gone.GetTypeName(t))
	}
}

func (s *engProvider) ProvideEngine(tagConf string) (Engine, error) {
	m, _ := gone.TagStringParse(tagConf)
	if v, ok := m[masterKey]; ok && (v == "" || cast.ToBool(v)) {
		group, err := s.xProvider.ProvideEngineGroup(tagConf)
		if err != nil {
			return nil, gone.ToError(err)
		}
		return newEng(group.Master(), s.logger), nil
	}
	if index, ok := m[slaveKey]; ok {
		i := cast.ToInt(index)
		group, err := s.xProvider.ProvideEngineGroup(tagConf)
		g.PanicIfErr(gone.ToErrorWithMsg(err, "can not create rocketmq consumer group"))
		slaves := group.Slaves()
		if i < 0 || i >= len(slaves) {
			return nil, gone.ToError("slave index out of range")
		}
		return newEng(slaves[i], s.logger), nil
	}

	provideEngine, err := s.xProvider.Provide(tagConf)
	if err != nil {
		return nil, gone.ToError(err)
	}
	return newEng(provideEngine, s.logger), nil
}
