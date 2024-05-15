package internal

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type GetterFunc = func(reflect.Value) (reflect.Value, bool)

func parseGetter(t reflect.Type, path string) (GetterFunc, reflect.Type, error) {
	ref := 0
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		ref++
	}
	switch t.Kind() {
	case reflect.Struct:
		if sf, ok := t.FieldByName(path); !ok {
			return nil, nil, fmt.Errorf("type %s no such field %s", t.String(), path)
		} else {
			return func(v reflect.Value) (reflect.Value, bool) {
				for i := 0; i < ref; i++ {
					v = v.Elem()
				}
				return v.FieldByIndex(sf.Index), true
			}, sf.Type, nil
		}
	case reflect.Map:
		switch t.Key().Kind() {
		case reflect.String:
			key := reflect.ValueOf(path)
			return func(v reflect.Value) (reflect.Value, bool) {
				for i := 0; i < ref; i++ {
					v = v.Elem()
				}
				return v.MapIndex(key), true
			}, t.Elem(), nil
		case reflect.Int:
			if k, err := strconv.Atoi(path); err != nil {
				return nil, nil, err
			} else {
				key := reflect.ValueOf(k)
				return func(v reflect.Value) (reflect.Value, bool) {
					for i := 0; i < ref; i++ {
						v = v.Elem()
					}
					return v.MapIndex(key), true
				}, t.Elem(), nil
			}
		case reflect.Int64:
			if k, err := strconv.ParseInt(path, 10, 64); err != nil {
				return nil, nil, err
			} else {
				key := reflect.ValueOf(k)
				return func(v reflect.Value) (reflect.Value, bool) {
					for i := 0; i < ref; i++ {
						v = v.Elem()
					}
					return v.MapIndex(key), true
				}, t.Elem(), nil
			}
		}
	case reflect.Slice, reflect.Array:
		if i, err := strconv.Atoi(path); err != nil {
			return nil, nil, err
		} else {
			return func(v reflect.Value) (reflect.Value, bool) {
				for i := 0; i < ref; i++ {
					v = v.Elem()
				}
				return v.Index(i), true
			}, t.Elem(), nil
		}
	}
	return nil, nil, fmt.Errorf("unsupported type and path")
}

func ParseGetters(t reflect.Type, path string) ([]GetterFunc, error) {
	paths := strings.Split(path, ".")
	rv := make([]GetterFunc, len(paths))
	for i, p := range paths {
		getter, nextType, err := parseGetter(t, p)
		if err != nil {
			return nil, err
		}
		rv[i] = getter
		t = nextType
	}
	return rv, nil
}
