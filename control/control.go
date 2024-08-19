package control

import (
	"context"
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
		timer = time.AfterFunc(after, f)
	}
}

// NewThrottle returns a throttled function that ensures `f` is not executed more frequently than `interval`.
func NewThrottle(interval time.Duration) func(f func()) bool {
	var (
		mu      sync.Mutex
		lastRun time.Time
	)

	return func(f func()) bool {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		if now.Sub(lastRun) >= interval {
			lastRun = now
			f()
			return true
		}

		return false
	}
}

// Interval calls `f` at regular `interval` until the `ctx` is cancelled.
func Interval(ctx context.Context, interval time.Duration, f func()) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				f()
			}
		}
	}()
}
