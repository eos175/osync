//go:build go1.23
// +build go1.23

package osync

import "iter"

// Iter iterates over current entries while holding a read lock.
// Iteration stops early when yield returns false.
// Map iteration order is not deterministic.
func (s *Map[K, T]) Iter() iter.Seq2[K, T] {
	return func(yield func(key K, value T) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		for k, v := range s.m {
			if !yield(k, v) {
				break
			}
		}
	}
}

// IterSnapshot iterates over a stable snapshot of current entries.
// Snapshot iteration does not hold locks while yielding.
// Snapshot order is not deterministic.
func (s *Map[K, T]) IterSnapshot() iter.Seq2[K, T] {
	snapshot := s.snapshot()
	return func(yield func(K, T) bool) {
		for _, e := range snapshot {
			if !yield(e.key, e.value) {
				break
			}
		}
	}
}

type tuple[K comparable, T any] struct {
	key   K
	value T
}

// Iter iterates over current keys while holding a read lock.
// Iteration stops early when yield returns false.
// Set iteration order is not deterministic.
func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(key T) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		for k := range s.m {
			if !yield(k) {
				break
			}
		}
	}
}

// IterSnapshot iterates over a stable snapshot of current keys.
// Snapshot iteration does not hold locks while yielding.
// Snapshot order is not deterministic.
func (s *Set[T]) IterSnapshot() iter.Seq[T] {
	snapshot := s.Keys()
	return func(yield func(T) bool) {
		for _, k := range snapshot {
			if !yield(k) {
				break
			}
		}
	}
}
