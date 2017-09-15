
[![GoDoc](https://godoc.org/github.com/tj/go-config?status.svg)](http://godoc.org/github.com/tj/go-config)
[![Build Status](https://travis-ci.org/tj/go-config.svg?branch=master)](https://travis-ci.org/tj/go-config)

# go-config

Simpler Go configuration with structs.

## Features

- Declare configuration with structs and tags
- Type coercion out of the box
- Validation out of the box
- Pluggable resolvers
- Built-in resolvers (flag, env)
- Unambiguous resolution (must be specified via `from`)

## Example

Source:

```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/tj/go-config"
)

type Options struct {
	Timeout     time.Duration `help:"message timeout"`
	Concurrency uint          `help:"message concurrency"`
	CacheSize   config.Bytes  `help:"cache size in bytes"`
	BatchSize   uint          `help:"batch size" validate:"min=1,max=1000"`
	LogLevel    string        `help:"log severity level" from:"env,flag"`
}

func main() {
	options := Options{
		Timeout:     time.Second * 5,
		Concurrency: 10,
		CacheSize:   config.ParseBytes("100mb"),
		BatchSize:   250,
	}

	config.MustResolve(&options)
	log.Printf("%+v", options)
}
```

Command-line:

```
$ LOG_LEVEL=error example --timeout 10s --concurrency 100 --cache-size 1gb
```

# License

MIT
