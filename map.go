package osync

import (
	"maps"
	"sync"
)

/*

https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c

*/

type Map[K comparable, T any] struct {
	mu sync.RWMutex
	m  map[K]T
}

func NewMap[K comparable, T any]() *Map[K, T] {
	return &Map[K, T]{m: make(map[K]T, minSizeSet)}
}

func (s *Map[K, T]) Get(key K) (T, bool) {
	s.mu.RLock()
	v, ok := s.m[key]
	s.mu.RUnlock()
	return v, ok
}

func (s *Map[K, T]) Set(key K, value T) {
	s.mu.Lock()
	s.m[key] = value
	s.mu.Unlock()
}

func (s *Map[K, T]) GetOrSet(key K, value T) (actual T, loaded bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.m[key]
	if ok {
		return val, false
	}
	s.m[key] = value
	return value, true
}

func (s *Map[K, T]) Clear() {
	s.mu.Lock()
	clear(s.m)
	s.mu.Unlock()
}

func (s *Map[K, T]) Delete(key K) {
	s.mu.Lock()
	delete(s.m, key)
	s.mu.Unlock()
}

func (s *Map[K, T]) Pop(key K) (T, bool) {
	s.mu.Lock()
	v, ok := s.m[key]
	if ok {
		delete(s.m, key)
	}
	s.mu.Unlock()
	return v, ok
}

func (s *Map[K, T]) ChangeKey(key, new_key K) (T, bool) {
	s.mu.Lock()
	v, ok := s.m[key]
	if ok {
		delete(s.m, key)
		s.m[new_key] = v
	}
	s.mu.Unlock()
	return v, ok
}

func (s *Map[K, T]) Len() int {
	s.mu.RLock()
	c := len(s.m)
	s.mu.RUnlock()
	return c
}

func (s *Map[K, T]) Clone() *Map[K, T] {
	s.mu.RLock()
	m := maps.Clone(s.m)
	s.mu.RUnlock()
	return &Map[K, T]{m: m}
}

func (s *Map[K, T]) ForEach(fn func(key K, value T) bool) {
	s.mu.RLock()
	for k, v := range s.m {
		if !fn(k, v) {
			break
		}
	}
	s.mu.RUnlock()
}
