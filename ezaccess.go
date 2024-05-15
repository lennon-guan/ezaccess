package ezaccess

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/lennon-guan/ezaccess/internal"
)

type PathStore struct {
	m sync.Map
}

type storeKey struct {
	t reflect.Type
	p string
}

func (s *PathStore) Load(t reflect.Type, path string) []internal.GetterFunc {
	if r, ok := s.m.Load(storeKey{t: t, p: path}); ok {
		return r.([]internal.GetterFunc)
	}
	return nil
}

func (s *PathStore) Store(t reflect.Type, path string, getters []internal.GetterFunc) {
	s.m.Store(storeKey{t: t, p: path}, getters)
}

func TryGet[T any](s *PathStore, m any, path string) (rv T, ok bool) {
	var (
		v       = reflect.ValueOf(m)
		vt      = v.Type()
		getters []internal.GetterFunc
		err     error
	)
	if s != nil {
		getters = s.Load(vt, path)
	}
	if getters == nil {
		getters, err = internal.ParseGetters(vt, path)
	}
	if err != nil {
		return
	}
	if s != nil {
		s.Store(vt, path, getters)
	}
	var _ok bool
	for _, getter := range getters {
		v, _ok = getter(v)
		if !_ok {
			return
		}
	}
	if v.IsValid() {
		rv, ok = v.Interface().(T)
	}
	return
}

func DefaultGet[T any](s *PathStore, m any, path string, def T) T {
	rv, ok := TryGet[T](s, m, path)
	if ok {
		return rv
	} else {
		return def
	}
}

func MustGet[T any](s *PathStore, m any, path string) T {
	rv, ok := TryGet[T](s, m, path)
	if ok {
		return rv
	}
	panic(fmt.Sprintf("get value by path %s failed", path))
}
