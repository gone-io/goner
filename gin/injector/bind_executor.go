package injector

import (
	"github.com/gone-io/gone/v2"
	"reflect"
)

var _ BindExecutor[any] = (*bindExecutor[any])(nil)

type bindExecutor[P any] struct {
	gone.Flag

	typeParsers []TypeParser[P] `gone:"*"`
	nameParsers []NameParser[P] `gone:"*"`

	typeParserMap map[reflect.Type]TypeParser[P]
	nameParserMap map[string]NameParser[P]
}

func (s *bindExecutor[P]) Init() error {
	s.typeParserMap = make(map[reflect.Type]TypeParser[P])
	for _, parser := range s.typeParsers {
		t := parser.Type()
		if _, ok := s.typeParserMap[t]; ok {
			return gone.NewInnerError("duplicate type parser", gone.InjectError)
		} else {
			s.typeParserMap[t] = parser
		}
	}

	s.nameParserMap = make(map[string]NameParser[P])
	for _, parser := range s.nameParsers {
		t := parser.Name()
		if _, ok := s.nameParserMap[t]; ok {
			return gone.NewInnerError("duplicate name parser", gone.InjectError)
		} else {
			s.nameParserMap[t] = parser
		}
	}
	return nil
}

func (s *bindExecutor[P]) InjectedByType(pt reflect.Type) BindFunc[P] {
	if parser, ok := s.typeParserMap[pt]; !ok {
		return nil
	} else {
		return func(context P) (reflect.Value, error) {
			return parser.Parse(context)
		}
	}
}

const anyName = "_*_"

func (s *bindExecutor[P]) FindFieldSetter(conf string, field reflect.StructField) (FieldSetter[P], error) {
	keyMap, keys := gone.TagStringParse(conf)
	if len(keys) > 0 && keys[0] != "" {
		name := keys[0]
		if keyMap[name] == "" {
			keyMap[anyName] = "true"
			keyMap[name] = field.Name
		}

		if nameParser, ok := s.nameParserMap[name]; !ok {
			return nil, gone.NewInnerError("can not find parser for field(name=%s)", gone.InjectError)
		} else {
			parser, err := nameParser.BuildParser(keyMap, field)
			if err != nil {
				return nil, err
			}
			return func(context P, fieldValue reflect.Value) error {
				if v, err := parser(context); err != nil {
					return err
				} else {
					fieldValue.Set(v)
					return nil
				}
			}, nil
		}
	}

	if bindFunc := s.InjectedByType(field.Type); bindFunc == nil {
		return nil, gone.NewInnerError("can not find parser for type %s", gone.InjectError)
	} else {
		return func(context P, fieldValue reflect.Value) error {
			if v, err := bindFunc(context); err != nil {
				return err
			} else {
				fieldValue.Set(v)
				return nil
			}
		}, nil
	}
}
