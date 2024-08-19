package osync

import (
	"context"
	"sync"
)

// Observable is a generic structure that represents a value that can be observed.
type Observable[T any] struct {
	value     T
	mu        sync.RWMutex
	observers []chan T
}

// NewObservable creates a new Observable with an initial value.
func NewObservable[T any](initialValue T) *Observable[T] {
	return &Observable[T]{
		value: initialValue,
	}
}

// Get returns the current value of the observable.
func (o *Observable[T]) Get() T {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.value
}

// Set updates the value of the observable and notifies all observers.
func (o *Observable[T]) Set(value T) {
	o.mu.Lock()
	o.value = value
	o.mu.Unlock()
	o.notifyObservers(value)
}

// Subscribe allows an observer to receive notifications when the value changes.
func (o *Observable[T]) Subscribe(ctx context.Context) <-chan T {
	ch := make(chan T, 1) // Buffered channel to avoid blocking.

	o.mu.Lock()
	o.observers = append(o.observers, ch)
	o.mu.Unlock()

	go func() {
		<-ctx.Done()
		o.removeObserver(ch)
		close(ch)
	}()

	return ch
}

// Len returns the number of observers currently subscribed.
func (o *Observable[T]) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.observers)
}

// notifyObservers sends the new value to all observers without blocking.
func (o *Observable[T]) notifyObservers(value T) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for _, observer := range o.observers {
		select {
		case observer <- value:
		default:
			// If the channel is full, skip sending to avoid blocking.
		}
	}
}

// removeObserver removes a channel from the list of observers by replacing it with the last channel.
func (o *Observable[T]) removeObserver(observer chan T) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for i, obs := range o.observers {
		if obs == observer {
			// Replace the removed channel with the last channel in the list.
			o.observers[i] = o.observers[len(o.observers)-1]
			o.observers = o.observers[:len(o.observers)-1]
			break
		}
	}
}
