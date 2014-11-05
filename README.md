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
    app.Get("/", func(resp husky.Response req *husky.Request) {
        // do something
    });
    ...

    // start server
    app.Serve()
}
```

That's all.


