# osync

**osync** is a Go library designed to provide thread-safe data structures and utilities for concurrent programming. Built with simplicity and performance in mind, `osync` leverages `sync.Mutex` and `sync.RWMutex` to ensure safe access to shared resources.

## Features

- **Thread-safe collections:** Protects against race conditions with minimal overhead.
- **Simple API:** Focuses on ease of use while offering powerful concurrency control.
- **Generic support:** Utilizes Go generics to create versatile and reusable data structures.

## Installation

To install `osync`, use `go get`:

```bash
go get github.com/eos175/osync
```


## Usage

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
}
```

## Documentation

The full documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/eos175/osync).


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

