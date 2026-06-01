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

// NewMap creates an empty concurrent map.
func NewMap[K comparable, T any]() *Map[K, T] {
	return &Map[K, T]{m: make(map[K]T, minSizeSet)}
}

// Get returns the value for key and whether it was present.
func (s *Map[K, T]) Get(key K) (T, bool) {
	s.mu.RLock()
	v, ok := s.m[key]
	s.mu.RUnlock()
	return v, ok
}

// Set stores value for key.
func (s *Map[K, T]) Set(key K, value T) {
	s.mu.Lock()
	s.m[key] = value
	s.mu.Unlock()
}

// GetOrSet returns the existing value for key if present.
// Otherwise it computes valueFn once under lock, stores it, and returns it.
func (s *Map[K, T]) GetOrSet(key K, valueFn func() T) (actual T, loaded bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.m[key]
	if ok {
		return val, true // Value was loaded
	}
	value := valueFn()
	s.m[key] = value
	return value, false // Value was set
}

// Clear removes all entries.
func (s *Map[K, T]) Clear() {
	s.mu.Lock()
	clear(s.m)
	s.mu.Unlock()
}

// Update atomically applies updateFn to the current value if key exists.
// It reports whether the key existed and was updated.
func (s *Map[K, T]) Update(key K, updateFn func(T) T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, exists := s.m[key]; exists {
		s.m[key] = updateFn(v)
		return true
	}
	return false
}

// UpdateIf updates the key if updateFn returns true alongside the new value.
// It reports whether an update was applied.
func (s *Map[K, T]) UpdateIf(key K, updateFn func(T) (T, bool)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, exists := s.m[key]; exists {
		if newVal, ok := updateFn(v); ok {
			s.m[key] = newVal
			return true
		}
	}
	return false
}

// Delete removes key if present.
func (s *Map[K, T]) Delete(key K) {
	s.mu.Lock()
	delete(s.m, key)
	s.mu.Unlock()
}

// DeleteIf removes key when it exists and condition(current) is true.
func (s *Map[K, T]) DeleteIf(key K, condition func(T) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, exists := s.m[key]; exists && condition(v) {
		delete(s.m, key)
	}
}

// Pop removes key and returns its value if present.
func (s *Map[K, T]) Pop(key K) (T, bool) {
	s.mu.Lock()
	v, ok := s.m[key]
	if ok {
		delete(s.m, key)
	}
	s.mu.Unlock()
	return v, ok
}

// PopIf removes key and returns its value when key exists and condition(current) is true.
func (s *Map[K, T]) PopIf(key K, condition func(T) bool) (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, exists := s.m[key]; exists && condition(v) {
		delete(s.m, key)
		return v, true
	}

	var zero T
	return zero, false
}

// RenameKey moves the value stored at key to newKey.
// It returns the previous value and whether key existed.
func (s *Map[K, T]) RenameKey(key, newKey K) (T, bool) {
	s.mu.Lock()
	v, ok := s.m[key]
	if ok {
		delete(s.m, key)
		s.m[newKey] = v
	}
	s.mu.Unlock()
	return v, ok
}

// Len returns the number of entries.
func (s *Map[K, T]) Len() int {
	s.mu.RLock()
	c := len(s.m)
	s.mu.RUnlock()
	return c
}

// Clone returns a shallow structural copy as a new concurrent Map.
func (s *Map[K, T]) Clone() *Map[K, T] {
	s.mu.RLock()
	m := maps.Clone(s.m)
	s.mu.RUnlock()
	return &Map[K, T]{m: m}
}

// Range iterates over current entries while holding a read lock.
// Iteration stops early when fn returns false.
// Map iteration order is not deterministic.
func (s *Map[K, T]) Range(fn func(key K, value T) bool) {
	s.mu.RLock()
	for k, v := range s.m {
		if !fn(k, v) {
			break
		}
	}
	s.mu.RUnlock()
}

// Snapshot returns a shallow snapshot of entries as key/value tuples.
// The returned slice can be iterated without holding locks.
// Snapshot order is not deterministic.
func (s *Map[K, T]) snapshot() []tuple[K, T] {
	s.mu.RLock()
	snapshot := make([]tuple[K, T], 0, len(s.m))
	for k, v := range s.m {
		snapshot = append(snapshot, tuple[K, T]{key: k, value: v})
	}
	s.mu.RUnlock()
	return snapshot
}
