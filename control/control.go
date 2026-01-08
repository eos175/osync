package control

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// NewDebouncer returns a debounced function that delays execution until after `after` duration has passed since the last call.
func NewDebouncer(after time.Duration) func(f func()) {
	var (
		mu    sync.Mutex
		timer *time.Timer
	)

	return func(f func()) {
		mu.Lock()
		defer mu.Unlock()

		if timer != nil {
			if !timer.Stop() {
				<-timer.C // Ensure the channel is drained
			}
		}
		timer = time.AfterFunc(after, func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("recovered from panic in NewDebouncer: %v\n", r)
				}
			}()
			f()
		})
	}
}

// NewThrottle returns a throttled function that ensures `f` is not executed more frequently than `interval`.
func NewThrottle(interval time.Duration) func(f func()) bool {
	var (
		mu      sync.Mutex
		lastRun time.Time
	)

	return func(f func()) (executed bool) {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		if now.Sub(lastRun) >= interval {
			lastRun = now
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("recovered from panic in NewThrottle: %v\n", r)
					executed = true
				}
			}()
			f()
			return true
		}

		return false
	}
}

// Interval calls `f` at regular `interval` until the `ctx` is cancelled.
// If `immediate` is true, `f` is called once at the beginning without waiting for the first tick.
func Interval(ctx context.Context, interval time.Duration, f func(), immediate ...bool) {
	go func() {
		safeF := func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("recovered from panic in Interval: %v\n", r)
				}
			}()
			f()
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		if len(immediate) > 0 && immediate[0] {
			safeF()
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				safeF()
			}
		}
	}()
}
