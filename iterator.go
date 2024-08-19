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
