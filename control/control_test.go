package control

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestIntervalPanicRecovery(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	// This function panics, which would normally crash the program
	f := func() {
		defer wg.Done()
		panic("intentional panic")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We wrap this in a subtest or just run it.
	// If Interval spawns a goroutine that panics, the test runner should fail/crash.
	// We use a small interval and cancel quickly.
	Interval(ctx, 10*time.Millisecond, f, true)

	// Wait for the function to run at least once
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Function ran
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for function execution")
	}
}

func TestDebouncerPanicRecovery(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	f := func() {
		defer wg.Done()
		panic("intentional panic in debouncer")
	}

	debounce := NewDebouncer(10 * time.Millisecond)
	debounce(f)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Function ran
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for debouncer execution")
	}
}

func TestThrottlePanicRecovery(t *testing.T) {
	// This function panics, which would normally crash the program
	f := func() {
		panic("intentional panic in throttle")
	}

	throttle := NewThrottle(10 * time.Millisecond)

	// Should not crash
	executed := throttle(f)
	if !executed {
		t.Error("expected throttled function to be executed")
	}
}
