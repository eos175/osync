package osync

import (
	"sync/atomic"
	"time"
	"unsafe"
)

/*

https://docs.python.org/3/library/asyncio-sync.html#asyncio.Event

https://gist.github.com/zviadm/c234426882bfc8acba88f3503edaaa36

https://gist.github.com/mkeeler/cb88cc762ca36733db0798ca80f1e73e

*/

// Event struct containing a state and a pointer to a channel
type Event struct {
	state   uint32
	channel unsafe.Pointer
}

// NewEvent initializes and returns a new Event instance
func NewEvent() *Event {
	ch := make(chan struct{})
	return &Event{
		channel: unsafe.Pointer(&ch),
	}
}

// broadcast closes the old channel and creates a new one atomically
func (e *Event) broadcast() {
	newCh := make(chan struct{})
	oldChPtr := atomic.SwapPointer(&e.channel, unsafe.Pointer(&newCh))
	close(*(*chan struct{})(oldChPtr))
}

// notifyChan returns a read-only channel for notification
func (e *Event) notifyChan() <-chan struct{} {
	chPtr := atomic.LoadPointer(&e.channel)
	return *(*chan struct{})(chPtr)
}

// IsSet checks if the event is set
func (e *Event) IsSet() bool {
	return atomic.LoadUint32(&e.state) == 1
}

// Set the event and broadcast if it's newly set
func (e *Event) Set() {
	if atomic.CompareAndSwapUint32(&e.state, 0, 1) {
		e.broadcast()
	}
}

// Clear the state of the event
func (e *Event) Clear() {
	atomic.StoreUint32(&e.state, 0)
}

// Wait until the event is set
func (e *Event) Wait() {
	for {
		if e.IsSet() {
			return
		}

		<-e.notifyChan()
	}

}

// https://github.com/golang/go/issues/9578

// WaitTimeout waits for the event to be set until the timeout
func (e *Event) WaitTimeout(timeout time.Duration) bool {
	if e.IsSet() { // Fast path
		return true
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		ch := e.notifyChan()

		select {
		case <-timer.C:
			return false
		case <-ch:
			if e.IsSet() {
				return true
			}
			// Spurious wakeup, continue loop
		}
	}
}
