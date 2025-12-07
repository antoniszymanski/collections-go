package typemap

import (
	"fmt"
	"iter"
	"maps"
	"reflect"
	"strings"
)

type Map struct {
	m map[reflect.Type]any
}

func New() Map {
	return Map{m: make(map[reflect.Type]any)}
}

func Get[T any](m Map) T {
	item, ok := m.m[reflect.TypeFor[T]()]
	if !ok {
		return *new(T)
	}
	return item.(T)
}

func Lookup[T any](m Map) (T, bool) {
	item, ok := m.m[reflect.TypeFor[T]()]
	if !ok {
		return *new(T), ok
	}
	return item.(T), ok
}

func Insert[T any](m Map, item T) {
	m.m[reflect.TypeFor[T]()] = item
}

func InsertAny(m Map, item any) {
	m.m[reflect.TypeOf(item)] = item
}

func Delete[T any](m Map) {
	delete(m.m, reflect.TypeFor[T]())
}

func Clear(m Map) {
	clear(m.m)
}

func Contains[T any](m Map) bool {
	_, exists := m.m[reflect.TypeFor[T]()]
	return exists
}

func All(m Map) iter.Seq2[reflect.Type, any] {
	return maps.All(m.m)
}

func Keys(m Map) iter.Seq[reflect.Type] {
	return maps.Keys(m.m)
}

func Values(m Map) iter.Seq[any] {
	return maps.Values(m.m)
}

func (m Map) String() string {
	var sb strings.Builder
	sb.WriteString("typemap.Map[")
	var i int
	for typ, val := range m.m {
		fmt.Fprintf(&sb, "%v:%v", typ, val)
		if i < len(m.m)-1 {
			sb.WriteByte(' ')
		}
		i++
	}
	sb.WriteByte(']')
	return sb.String()
}
