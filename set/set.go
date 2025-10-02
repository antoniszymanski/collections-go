// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"encoding/json"
	"fmt"
	"iter"
	"maps"
)

type Set[T comparable] struct {
	set map[T]bool
}

func New[T comparable](size int) Set[T] {
	return Set[T]{set: make(map[T]bool, max(0, size))}
}

func From[T comparable](items ...T) Set[T] {
	s := New[T](len(items))
	s.InsertMany(items...)
	return s
}

func Collect[T comparable](seq iter.Seq[T]) Set[T] {
	s := New[T](0)
	s.InsertSeq(seq)
	return s
}

func (s Set[T]) Insert(item T) {
	s.set[item] = true
}

func (s Set[T]) InsertMany(items ...T) {
	for _, item := range items {
		s.Insert(item)
	}
}

func (s Set[T]) InsertSeq(seq iter.Seq[T]) {
	for item := range seq {
		s.Insert(item)
	}
}

func (s Set[T]) Delete(item T) {
	delete(s.set, item)
}

func (s Set[T]) Clear() {
	clear(s.set)
}

func (s Set[T]) Contains(item T) bool {
	return s.set[item]
}

func (s Set[T]) Size() int {
	return len(s.set)
}

func (s Set[T]) Empty() bool {
	return len(s.set) == 0
}

func (s Set[T]) Equal(other Set[T]) bool {
	if len(s.set) != len(other.set) {
		return false
	}
	for item := range s.set {
		if !other.set[item] {
			return false
		}
	}
	return true
}

func (s Set[T]) Clone() Set[T] {
	return Set[T]{set: maps.Clone(s.set)}
}

func (s Set[T]) Items() []T {
	items := make([]T, 0, len(s.set))
	for item := range s.set {
		items = append(items, item)
	}
	return items
}

func (s Set[T]) All() iter.Seq[T] {
	return maps.Keys(s.set)
}

func (s Set[T]) String() string {
	return fmt.Sprint(s.Items())
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Items())
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	if s.set == nil {
		s.set = make(map[T]bool, len(items))
	}
	s.InsertMany(items...)
	return nil
}
