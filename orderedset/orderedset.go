// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package orderedset

import (
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"slices"
)

type OrderedSet[T comparable] struct {
	set   map[T]bool
	items []T
}

func New[T comparable](size int) *OrderedSet[T] {
	return &OrderedSet[T]{
		set:   make(map[T]bool, max(0, size)),
		items: make([]T, 0, max(0, size)),
	}
}

func From[T comparable](items ...T) *OrderedSet[T] {
	s := New[T](len(items))
	s.InsertMany(items...)
	return s
}

func Collect[T comparable](seq iter.Seq[T]) *OrderedSet[T] {
	s := New[T](0)
	s.InsertSeq(seq)
	return s
}

func (s *OrderedSet[T]) Insert(item T) (modified bool) {
	if modified = !s.set[item]; modified {
		s.set[item] = true
		s.items = append(s.items, item)
	}
	return
}

func (s *OrderedSet[T]) InsertMany(items ...T) (modified bool) {
	for _, item := range items {
		if s.Insert(item) {
			modified = true
		}
	}
	return
}

func (s *OrderedSet[T]) InsertSeq(seq iter.Seq[T]) (modified bool) {
	for item := range seq {
		if s.Insert(item) {
			modified = true
		}
	}
	return
}

func (s *OrderedSet[T]) Delete(item T) {
	if s == nil {
		return
	}
	delete(s.set, item)
	if i := slices.Index(s.items, item); i != -1 {
		s.items = slices.Delete(s.items, i, i)
	}
}

func (s *OrderedSet[T]) Clear() {
	if s != nil {
		clear(s.set)
		clear(s.items)
	}
}

func (s *OrderedSet[T]) Contains(item T) bool {
	return s != nil && s.set[item]
}

func (s *OrderedSet[T]) Size() int {
	if s == nil {
		return 0
	}
	return len(s.items)
}

func (s *OrderedSet[T]) Empty() bool {
	return s == nil || len(s.items) == 0
}

func (s *OrderedSet[T]) Equal(other *OrderedSet[T]) bool {
	if s.Empty() && other.Empty() {
		return true
	} else if len(s.items) != len(other.items) {
		return false
	}
	for i := range s.items {
		if s.items[i] != other.items[i] {
			return false
		}
	}
	return true
}

func (s *OrderedSet[T]) Clone() *OrderedSet[T] {
	if s == nil {
		return nil
	}
	return &OrderedSet[T]{
		set:   maps.Clone(s.set),
		items: slices.Clone(s.items),
	}
}

func (s *OrderedSet[T]) Items() []T {
	if s == nil {
		return nil
	}
	return slices.Clone(s.items)
}

func (s *OrderedSet[T]) All() iter.Seq2[int, T] {
	if s == nil {
		return func(yield func(int, T) bool) {}
	}
	return slices.All(s.items)
}

func (s *OrderedSet[T]) String() string {
	return fmt.Sprint(s.items)
}

func (s OrderedSet[T]) MarshalJSON() ([]byte, error) {
	if s.Empty() {
		return []byte("[]"), nil
	}
	return json.Marshal(s.items)
}

func (s *OrderedSet[T]) UnmarshalJSON(data []byte) error {
	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	if s.set == nil {
		s.set = make(map[T]bool, len(items))
	}
	if cap(s.items) == 0 {
		for _, item := range items {
			s.set[item] = true
		}
		s.items = items
	} else {
		s.InsertMany(items...)
	}
	return nil
}
