# osync

**osync** is a Go library designed to provide thread-safe data structures and utilities for concurrent programming. Built with simplicity and performance in mind, `osync` leverages `Mutex`, `RWMutex` and `Atomic` to ensure safe access to shared resources.

## Features

- **Thread-safe collections:** Protects against race conditions with minimal overhead.
- **Observable values:** Allows observing changes to a value.
- **Event handling:** Provides synchronization primitives for coordinating tasks.
- **Simple API:** Focuses on ease of use while offering powerful concurrency control.
- **Generic support:** Utilizes Go generics to create versatile and reusable data structures.

## Installation

To install `osync`, use `go get`:

```bash
go get github.com/eos175/osync
```

## Usage

### Set Example

Here's an example of how to use the `Set` provided by `osync`:

```go
package main

import (
	"fmt"

	"github.com/eos175/osync"
)

func main() {
	set := osync.NewSet[int]()

	set.Add(1)
	set.Add(2)
	set.Add(3)

	fmt.Println("Set has 2:", set.Has(2)) // Output: Set has 2: true
	fmt.Println("Set length:", set.Len()) // Output: Set length: 3

	set.Delete(2)

	fmt.Println("Set has 2:", set.Has(2)) // Output: Set has 2: false

	// go1.23
	fmt.Println("Set contents:")
	for key := range set.Iterator() {
		fmt.Println(key)
	}
}
```

### Observable Example

Here's an example of how to use the `Observable` provided by `osync`:

```go
package main

import (
	"fmt"
	"time"

	"github.com/eos175/osync"
)

func main() {
	obs := osync.NewObservable[int](0)
	defer obs.Close()

	// Subscribe to changes. This returns a channel and an unsubscribe function.
	ch, unsubscribe := obs.Subscribe()
	// It's important to call unsubscribe when done to avoid leaks.
	defer unsubscribe()

	// This goroutine will stop the subscription after 5 seconds.
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Unsubscribing...")
		unsubscribe()
	}()

	// This goroutine updates the observable value.
	go func() {
		for i := 1; ; i++ {
			// This Set will be missed if it happens after unsubscribe.
			obs.Set(i * i)
			time.Sleep(1 * time.Second)
		}
	}()

	// Print updates received from the observable.
	// The loop will end when the channel is closed by the unsubscribe call.
	for value := range ch {
		fmt.Println("Received value:", value)
	}

	fmt.Println("Subscription ended.")
}
```

### Event Example

Here's an example of how to use the `Event` provided by `osync`:

```go
package main

import (
	"fmt"
	"time"

	"github.com/eos175/osync"
)

func main() {
	event := osync.NewEvent()

	go func() {
		// Wait for the event to be set
		fmt.Println("Waiting for event to be set...")
		event.Wait()
		fmt.Println("Event is set!")
	}()

	go func() {
		// Simulate some work before setting the event
		time.Sleep(2 * time.Second)
		fmt.Println("Setting event...")
		event.Set()
	}()

	// Wait for the event to be set
	time.Sleep(3 * time.Second)
}
```

### Control Utilities

The `control` package provides utilities for controlling function execution flow.

#### Debouncer

Executes a function only after a specified duration has passed without new calls.

```go
package main

import (
	"fmt"
	"time"

	"github.com/eos175/osync/control"
)

func main() {
	debouncer := control.NewDebouncer(100 * time.Millisecond)

	// Will execute only the last call after 100ms
	debouncer(func() { fmt.Println("1") })
	debouncer(func() { fmt.Println("2") })
	debouncer(func() { fmt.Println("3") }) // Only this one runs

	time.Sleep(200 * time.Millisecond)
}
```

#### Throttle

Ensures a function is not executed more frequently than a specified interval.

```go
package main

import (
	"fmt"
	"time"

	"github.com/eos175/osync/control"
)

func main() {
	throttle := control.NewThrottle(100 * time.Millisecond)

	// First call runs immediately
	throttle(func() { fmt.Println("Run 1") })

	// This call is skipped because it's too soon
	throttle(func() { fmt.Println("Run 2") })

	time.Sleep(150 * time.Millisecond)
	// This call runs
	throttle(func() { fmt.Println("Run 3") })
}
```

#### Scheduled Periodic Tasks

To schedule a task to run periodically starting at a specific time (e.g., daily backups), you can combine `NextDailyAt` with `IntervalAt`.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/eos175/osync/control"
)

func main() {
	ctx := context.Background()

	// Calculate the next occurrence of 03:00:00 AM (local time)
	start := control.NextDailyAt(time.Now(), 3, 0, 0)

	// To use a specific timezone (e.g., UTC):
	// start := control.NextDailyAt(time.Now().In(time.UTC), 3, 0, 0)

	// Schedule the task to run daily starting at 'start'
	control.IntervalAt(ctx, start, 24*time.Hour, func() {
		fmt.Println("Starting daily backup...")
	})

	// Block main process
	select {}
}
```

#### IntervalAt

Execute a task at a regular interval, starting at a specific future time.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/eos175/osync/control"
)

func main() {
	ctx := context.Background()
	start := time.Now().Add(1 * time.Hour) // Start in 1 hour

	control.IntervalAt(ctx, start, 30*time.Minute, func() {
		fmt.Println("Running task...")
	})
}
```

## Documentation

The full documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/eos175/osync).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
