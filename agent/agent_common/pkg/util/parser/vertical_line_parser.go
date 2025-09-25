package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type verticalLineFieldSetter func(v reflect.Value, raw string) error

type VerticalLineParser[T any] struct {
	fieldMap map[string]verticalLineFieldSetter
	fieldIdx map[string]int

	sep string
}

func (p *VerticalLineParser[T]) makeSetter(kind string) verticalLineFieldSetter {
	return func(field reflect.Value, raw string) error {
		raw = strings.TrimSpace(raw)
		switch kind {
		case "string":
			field.SetString(raw)
		case "int":
			val, err := strconv.Atoi(raw)
			if err != nil {
				return err
			}
			field.SetInt(int64(val))
		case "float64":
			if raw == "NaN" {
				field.SetFloat(0)
				return nil
			}
			val, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return err
			}
			field.SetFloat(val)
		default:
			return fmt.Errorf("unsupported kind: %s", kind)
		}
		return nil
	}
}

func CreateVerticalLineParser[T any](sep string) *VerticalLineParser[T] {
	var t T
	typ := reflect.TypeOf(t)
	if typ.Kind() != reflect.Struct {
		panic("T must be a struct")
	}

	parser := &VerticalLineParser[T]{
		fieldMap: make(map[string]verticalLineFieldSetter),
		fieldIdx: make(map[string]int),
		sep : sep,
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("agent_common_parser")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		key := strings.TrimSpace(parts[0])
		kind := "string"
		if len(parts) > 1 {
			kind = parts[1]
		}

		parser.fieldIdx[key] = i
		parser.fieldMap[key] = parser.makeSetter(kind)
	}

	return parser
}

func (p *VerticalLineParser[T]) Load(input string, output *T) error {
	lines := strings.Split(input, "\n")
	val := reflect.ValueOf(output).Elem()

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, p.sep, 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		setter, ok := p.fieldMap[key]
		if !ok {
			continue
		}

		idx := p.fieldIdx[key]
		field := val.Field(idx)

		err := setter(field, value)
		if err != nil {
			return fmt.Errorf("error parsing key %q with value %q: %w", key, value, err)
		}
	}

	return nil
}