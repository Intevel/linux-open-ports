# linux-open-ports

📡 Go package that retrieves information about open network ports on a Linux system. It identifies active ports and the processes associated with them by reading the system's network connection files in `/proc/net/`.

## Features

- Fetches open TCP and UDP ports.
- Retrieves the PID associated with each open port.
- Supports both IPv4 and IPv6 connections.
- Efficiently handles connections by avoiding duplicates in the port list.

## Installation

To use the `linuxopenports` package, you can add it to your Go project by importing it:

```go
import "github.com/intevel/linux-open-ports"

```

## Usage

The `linuxopenports` package provides a single function, `GetOpenPorts()`, which returns a list of open ports and their associated processes. The function signature is as follows:

```go

func GetOpenPorts() ([]linuxopenports.Port, error)

```

The `GetOpenPorts()` function returns a slice of `Port` structs, each containing the following fields:

```go

type Port struct {
    Protocol string // The protocol used by the port (TCP or UDP).
    Port     uint16 // The port number.
    PID      int    // The process ID associated with the port.
    Program  string // The name of the program associated with the port.
}

```

Here is an example of how to use the `GetOpenPorts()` function:

```go

package main

import (
    "fmt"
    "github.com/intevel/linux-open-ports"
)

func main() {
    openPorts, err := linuxopenports.GetOpenPorts()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    for _, port := range openPorts {
        fmt.Printf("Protocol: %s, Port: %d, PID: %d, Program: %s\n", port.Protocol, port.Port, port.PID, port.Program)
    }
}

```

## License

Published under MIT - Made with ❤️ by Conner Bachmann