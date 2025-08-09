//go:build go1.23
// +build go1.23

package osync

import "iter"

func (s *Map[K, T]) Iterator() iter.Seq2[K, T] {
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

// Create a slice of key-value pairs
type Tuple[K comparable, T any] struct {
	key K
	val T
}

// SnapshotIterator returns an iterator over an immutable copy (snapshot) of the map.
// This allows safe iteration without blocking concurrent access to the original map.
func (s *Map[K, T]) SnapshotIterator() iter.Seq2[K, T] {
	// Take a snapshot under read lock
	s.mu.RLock()
	snapshot := make([]Tuple[K, T], 0, len(s.m))
	for k, v := range s.m {
		snapshot = append(snapshot, Tuple[K, T]{key: k, val: v})
	}
	s.mu.RUnlock()

	// Return an iterator over the snapshot
	return func(yield func(K, T) bool) {
		for _, e := range snapshot {
			if !yield(e.key, e.val) {
				break
			}
		}
	}
}

func (s *Set[T]) Iterator() iter.Seq[T] {
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

func (s *Set[T]) SnapshotIterator() iter.Seq[T] {
	s.mu.RLock()
	// Take snapshot of keys
	snapshot := make([]T, 0, len(s.m))
	for k := range s.m {
		snapshot = append(snapshot, k)
	}
	s.mu.RUnlock()

	// Return iterator over snapshot
	return func(yield func(T) bool) {
		for _, k := range snapshot {
			if !yield(k) {
				break
			}
		}
	}
}
