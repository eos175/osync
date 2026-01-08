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
	first := interval
	if len(immediate) > 0 && immediate[0] {
		first = 0
	}

	runTicker(ctx, first, interval, f)
}

// IntervalAt calls `f` starting at `start` time and then every `interval` until `ctx` is cancelled.
// If `start` is in the past, `f` is executed immediately.
func IntervalAt(ctx context.Context, start time.Time, interval time.Duration, f func()) {
	delay := time.Until(start)
	if delay < 0 {
		delay = 0
	}

	runTicker(ctx, delay, interval, f)
}

func runTicker(
	ctx context.Context,
	firstDelay time.Duration,
	interval time.Duration,
	f func(),
) {
	go func() {
		safeF := func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("recovered from panic: %v\n", r)
				}
			}()
			f()
		}

		if firstDelay > 0 {
			timer := time.NewTimer(firstDelay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		}

		safeF()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

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
