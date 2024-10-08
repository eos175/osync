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
	"context"
	"fmt"
	"time"

	"github.com/eos175/osync"
)

func main() {
	obs := osync.NewObservable[int](0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to changes
	ch := obs.Subscribe(ctx)

	go func() {
		// Update the observable value
		for i := 0; i < 10; i++ {
			obs.Set(i*i)
			time.Sleep(1 * time.Second)
		}
	}()

	// Print updates received from the observable
	for value := range ch {
		fmt.Println("Received value:", value)
	}
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


## Documentation

The full documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/eos175/osync).


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

