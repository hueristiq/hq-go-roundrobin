# hqgoroundrobin

[![go report card](https://goreportcard.com/badge/github.com/hueristiq/hqgoroundrobin)](https://goreportcard.com/report/github.com/hueristiq/hqgoroundrobin) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hqgoroundrobin.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hqgoroundrobin/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hqgoroundrobin.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hqgoroundrobin/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hqgoroundrobin/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hqgoroundrobin/blob/master/CONTRIBUTING.md)

A [Go(Golang)](http://golang.org/) high-quality, concurrency-safe implementation of [Round Robin(RR)](https://en.wikipedia.org/wiki/Round-robin_scheduling) algorithm for managing and cycling through a collection of items. This package is designed to be easy to integrate and use for load balancing, task distribution, and other scenarios where items need to be processed or served in a cyclic order.

## Resources

* [Features](#features)
* [Installation](#installation)
* [Usage](#usage)
* [Contributing](#contributing)
* [Licensing](#licensing)
* [Credits](#credits)
    * [Contributors](#contributors)
    * [Similar Projects](#similar-projects)

## Features

* Ensures safe concurrent access and modification of the round-robin queue.
* Prevents duplicate items in the queue, maintaining the integrity of the rotation.
* Customizable configuration to define how often the rotation should move to the next item.
* Provides a straightforward API for adding items and retrieving the next item in the round-robin sequence.

## Installation

To install `hqgoroundrobin`, use the go get command:

```bash
go get -v -u github.com/hueristiq/hqgoroundrobin
```

## Usage

Here's a simple example to get you started with `hqgoroundrobin`:

```go
package main

import (
	"fmt"
	
	roundrobin "github.com/hueristiq/hqgoroundrobin"
)

func main() {
	// Create a new round-robin instance with default options and initial items.
	rr, err := roundrobin.New("item1", "item2", "item3")
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

[Issues](https://github.com/hueristiq/hqgoroundrobin/issues) and [Pull Requests](https://github.com/hueristiq/hqgoroundrobin/pulls) are welcome! **Check out the [contribution guidelines](https://github.com/hueristiq/hqgoroundrobin/blob/master/CONTRIBUTING.md).**

## Licensing

This utility is distributed under the [MIT license](https://github.com/hueristiq/hqgoroundrobin/blob/master/LICENSE).

## Credits

### Contributors

Thanks to the amazing [contributors](https://github.com/hueristiq/hqgoroundrobin/graphs/contributors) for keeping this project alive.

[![contributors](https://contrib.rocks/image?repo=hueristiq/hqgoroundrobin&max=500)](https://github.com/hueristiq/hqgoroundrobin/graphs/contributors)

### Similar Projects

Thanks to similar open source projects - check them out, may fit in your needs

[hlts2/round-robin](https://github.com/hlts2/round-robin) ◇ [projectdiscovery/roundrobin](https://github.com/projectdiscovery/roundrobin)