# hq-go-roundrobin

![made with go](https://img.shields.io/badge/made%20with-Go-1E90FF.svg) [![go report card](https://goreportcard.com/badge/github.com/hueristiq/hq-go-roundrobin)](https://goreportcard.com/report/github.com/hueristiq/hq-go-roundrobin) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hq-go-roundrobin/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hq-go-roundrobin.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-roundrobin/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hq-go-roundrobin.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-roundrobin/issues?q=is:issue+is:closed) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hq-go-roundrobin/blob/master/CONTRIBUTING.md)

`hq-go-roundrobin` is a [Go (Golang)](http://golang.org/) high-quality, concurrency-safe implementation of [Round Robin(RR)](https://en.wikipedia.org/wiki/Round-robin_scheduling) algorithm for managing and cycling through a collection of items.

## Resources

* [Features](#features)
* [Installation](#installation)
* [Usage](#usage)
* [Contributing](#contributing)
* [Licensing](#licensing)

## Features

* Ensures safe concurrent access and modification of the round-robin queue.
* Prevents duplicate items in the queue, maintaining the integrity of the rotation.
* Customizable configuration to define how often the rotation should move to the next item.
* Provides a straightforward API for adding items and retrieving the next item in the round-robin sequence.

## Installation

To install `hq-go-roundrobin`, run:

```bash
go get -v -u github.com/hueristiq/hq-go-roundrobin
```

Make sure your Go environment is set up properly (Go 1.x or later is recommended).

## Usage

Here's a simple example to get you started with `hq-go-roundrobin`:

```go
package main

import (
	"fmt"
	
	hqgoroundrobin "github.com/hueristiq/hq-go-roundrobin"
)

func main() {
	// Create a new round-robin instance with default options and initial items.
	rr, err := hqgoroundrobin.New("item1", "item2", "item3")
	if err != nil {
		panic(err)
	}

	// Add more items if needed
	rr.Add("item4", "item5")

	// Retrieve and process items in a round-robin fashion
	for i := 0; i < 10; i++ {
		item := rr.Next()
		fmt.Printf("Serving: %s\n", item.Value())
	}

	/*
		Output:
		item1
		item2
		item3
		item4
		item5
		item1
		item2
		item3
		item4
		item5
	*/

}
```

## Contributing

Contributions are welcome and encouraged! Feel free to submit [Pull Requests](https://github.com/hueristiq/hq-go-roundrobin/pulls) or report [Issues](https://github.com/hueristiq/hq-go-roundrobin/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hq-go-roundrobin/blob/master/CONTRIBUTING.md).

A big thank you to all the [contributors](https://github.com/hueristiq/hq-go-roundrobin/graphs/contributors) for your ongoing support!

![contributors](https://contrib.rocks/image?repo=hueristiq/hq-go-roundrobin&max=500)

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hq-go-roundrobin/blob/master/LICENSE).