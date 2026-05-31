package osync

import (
	"slices"
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

// Set updates the current value and notifies current observers.
// Notification is best-effort: if a subscriber channel buffer is full, that
// subscriber does not receive this update.
func (o *Observable[T]) Set(value T) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.value = value

	for _, observer := range o.observers {
		select {
		case observer <- value:
		default:
			// If the channel is full, skip sending to avoid blocking.
		}
	}
}

// Subscribe registers a new observer.
// It returns a channel for updates and an unsubscribe function.
// The current value is sent to the subscriber upon subscription.
func (o *Observable[T]) Subscribe() (<-chan T, func()) {
	ch := make(chan T, 1)

	o.mu.Lock()
	o.observers = append(o.observers, ch)
	// Send initial value while holding the lock to ensure atomicity.
	// This is safe and won't deadlock because the channel is buffered and new.
	ch <- o.value
	o.mu.Unlock()

	unsubscribe := func() {
		// If we successfully remove the channel, this call owns closing it.
		if o.removeObserver(ch) {
			close(ch)
		}
	}
	return ch, unsubscribe
}

// Len returns the number of observers currently subscribed.
func (o *Observable[T]) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.observers)
}

// UnsubscribeAll terminates all current subscriptions.
// The observable remains active and can accept new subscriptions.
func (o *Observable[T]) UnsubscribeAll() {
	o.mu.Lock()
	// Atomically claim all remaining observers.
	observersToClose := o.observers
	o.observers = nil
	o.mu.Unlock()

	// And close them.
	for _, ch := range observersToClose {
		close(ch)
	}
}

// removeObserver removes a channel from the list and returns true if it was found.
// This atomicity is key to deciding which routine is responsible for closing the channel.
func (o *Observable[T]) removeObserver(observer chan T) bool {
	o.mu.Lock()
	defer o.mu.Unlock()

	if index := slices.Index(o.observers, observer); index != -1 {
		// Replace the removed channel with the last channel in the list.
		o.observers[index] = o.observers[len(o.observers)-1]
		o.observers = o.observers[:len(o.observers)-1]
		return true // Success
	}
	return false // Not found
}
