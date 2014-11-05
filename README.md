Husky
=====

A Simple Router, HTTP, utility for golang

## Installation

Use `go get` command:

```
go get github.com/ysugimoto/husky
```

## Usage

Import and create app, and bind routes like Sinatra:

```
package main

import "github.com/ysugimoto/husky"


func main() {
    // create app
    app := husky.NewApp()

    // bind routes as you need
    app.Get("/", func(disp *husky.Dispatcher) {
        disp.Output.SetStatus(200)
        disp.Output.SetHeader("Content-Type", "text/plain")
        disp.Output.Send("Hello, Husky!")
    });
    ...

    // start server
    app.Serve()
}
```

That's all.

### Configuration

Make `config.yml` at current directory (e.g main.go) to override config and write:

```
# listen address
host: 127.0.0.1

# listen port
port: 8888

# listen path
path: /
```

Will listen `127.0.0.1:8888`.
