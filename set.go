package osync

import (
	"maps"
	"sync"
)

const minSizeSet = 64

type null = struct{}

type Set[T comparable] struct {
	mu sync.RWMutex
	m  map[T]null
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]null, minSizeSet)}
}

// Removes all keys from the set.
func (s *Set[T]) Clear() {
	s.mu.Lock()
	clear(s.m)
	s.mu.Unlock()
}

// Adds a key to the set. Returns `true` if the key was added, or `false` if it already existed.
func (s *Set[T]) Add(key T) bool {
	s.mu.Lock()
	prevLen := len(s.m)
	s.m[key] = null{}
	cLen := len(s.m)
	s.mu.Unlock()
	return prevLen != cLen
}

// Checks if a key exists in the set.
func (s *Set[T]) Has(key T) bool {
	s.mu.RLock()
	_, ok := s.m[key]
	s.mu.RUnlock()
	return ok
}

// Removes a key from the set.
func (s *Set[T]) Delete(key T) {
	s.mu.Lock()
	delete(s.m, key)
	s.mu.Unlock()
}

// Removes a key from the set and returns `true` if the key existed.
func (s *Set[T]) Pop(key T) bool {
	s.mu.Lock()
	_, ok := s.m[key]
	if ok {
		delete(s.m, key)
	}
	s.mu.Unlock()
	return ok
}

// Iterates over all keys in the set, applying the provided function.
func (s *Set[T]) ForEach(fn func(key T) bool) {
	s.mu.RLock()
	for key := range s.m {
		if !fn(key) {
			break
		}
	}
	s.mu.RUnlock()
}

// Returns a slice of all keys in the set.
func (s *Set[T]) Keys() []T {
	s.mu.RLock()
	keys := make([]T, 0, len(s.m))
	for key := range s.m {
		keys = append(keys, key)
	}
	s.mu.RUnlock()
	return keys
}

// Returns the number of keys in the set.
func (s *Set[T]) Len() int {
	s.mu.RLock()
	c := len(s.m)
	s.mu.RUnlock()
	return c
}

// Returns a shallow copy of the set.
func (s *Set[T]) Clone() *Set[T] {
	s.mu.RLock()
	m := maps.Clone(s.m)
	s.mu.RUnlock()
	return &Set[T]{m: m}
}
