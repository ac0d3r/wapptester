package wapptester

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/PaesslerAG/gval"
)

var ErrUnsupportedType = errors.New("unsupported types")

type SelectWrapper struct {
	attrs map[string]any
}

var _ gval.Selector = &SelectWrapper{}

func NewSelectWrapper(v any, tag string) (gval.Selector, error) {
	var isptr bool

	vt := reflect.TypeOf(v)
	vtkind := vt.Kind()
	if vtkind == reflect.Ptr {
		isptr = true
		vtkind = vt.Elem().Kind()
	}
	if vtkind != reflect.Struct {
		return nil, ErrUnsupportedType
	}

	attrs := make(map[string]any)
	// fileds
	vv := reflect.ValueOf(v)
	if isptr {
		vv = vv.Elem()
	}
	for i := 0; i < vv.NumField(); i++ {
		t := vv.Type().Field(i).Tag.Get(tag)
		if t != "" {
			attrs[t] = vv.Field(i).Interface()
		}
	}

	// methods
	for i := 0; i < vv.NumMethod(); i++ {
		attrs[vv.Type().Method(i).Name] = vv.Method(i).Interface()
	}

	return &SelectWrapper{attrs: attrs}, nil
}

func (s *SelectWrapper) SelectGVal(c context.Context, key string) (interface{}, error) {
	v, ok := s.attrs[key]
	if !ok {
		return nil, fmt.Errorf("unknown parameter %s", key)
	}
	return v, nil
}
