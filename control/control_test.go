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

func TestIntervalAt(t *testing.T) {
	var mu sync.Mutex
	executionTimes := []time.Time{}

	f := func() {
		mu.Lock()
		defer mu.Unlock()
		executionTimes = append(executionTimes, time.Now())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start 100ms in the future
	startDelay := 100 * time.Millisecond
	interval := 50 * time.Millisecond
	startTime := time.Now().Add(startDelay)

	IntervalAt(ctx, startTime, interval, f)

	// Wait enough time for: initial delay (100) + one interval (50) + buffer
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(executionTimes) < 2 {
		t.Fatalf("expected at least 2 executions, got %d", len(executionTimes))
	}

	// Verify first execution was close to startTime
	firstExec := executionTimes[0]
	diff := firstExec.Sub(startTime)
	if diff < -20*time.Millisecond || diff > 20*time.Millisecond {
		t.Errorf("first execution time off. Expected around %v, got %v (diff: %v)", startTime, firstExec, diff)
	}

	// Verify second execution was interval later
	secondExec := executionTimes[1]
	expectedSecond := firstExec.Add(interval)
	diff2 := secondExec.Sub(expectedSecond)
	if diff2 < -20*time.Millisecond || diff2 > 20*time.Millisecond {
		t.Errorf("second execution time off. Expected around %v, got %v (diff: %v)", expectedSecond, secondExec, diff2)
	}
}
