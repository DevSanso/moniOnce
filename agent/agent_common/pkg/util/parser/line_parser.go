package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type linePaserFieldSetter func(v reflect.Value, raw string) error

type LinePaser[T any] struct {
	fieldMap map[int]linePaserFieldSetter
	fieldFn func(string)[]string
}

func (l *LinePaser[T]) makeSetter(kind string) linePaserFieldSetter {
	return func(field reflect.Value, raw string) error {
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

func CreateLinePaser[T any](sep string) *LinePaser[T] {
	var t T
	typ := reflect.TypeOf(t)
	if typ.Kind() != reflect.Struct {
		panic("T must be a struct")
	}

	loader := &LinePaser[T]{
		fieldMap: make(map[int]linePaserFieldSetter),
	}

	if sep == " " {
		loader.fieldFn = strings.Fields	
	} else {
		loader.fieldFn = func(line string) []string {
			return strings.Split(line, sep)
		}
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("agent_common_parser")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		if len(parts) != 2 {
			panic(fmt.Sprintf("invalid tag format for field %s", field.Name))
		}

		pos, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Sprintf("invalid position in tag for field %s", field.Name))
		}

		kind := parts[1]
		setter := loader.makeSetter(kind)

		idx := i
		loader.fieldMap[pos] = func(v reflect.Value, raw string) error {
			f := v.Field(idx)
			return setter(f, raw)
		}
	}

	return loader
}

func (l *LinePaser[T]) Load(line string, output *T) error {
	tokens := l.fieldFn(line)
	val := reflect.ValueOf(output).Elem()

	for pos, setter := range l.fieldMap {
		if pos >= len(tokens) {
			return fmt.Errorf("not enough fields for position %d", pos)
		}
		err := setter(val, tokens[pos])
		if err != nil {
			return fmt.Errorf("error setting field at pos %d: %w", pos, err)
		}
	}
	return nil
}