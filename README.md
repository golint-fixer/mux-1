# vinci/mux [![Build Status](https://travis-ci.org/vinci-proxy/gentleman.png)](https://travis-ci.org/vinci-proxy/gentleman) [![GoDoc](https://godoc.org/github.com/vinci-proxy/vinci/mux?status.svg)](https://godoc.org/github.com/vinci-proxy/vinci/mux) [![API](https://img.shields.io/badge/status-stable-green.svg?style=flat)](https://godoc.org/github.com/vinci-proxy/vinci/mux) [![Go Report Card](https://goreportcard.com/badge/github.com/vinci-proxy/vinci/mux)](https://goreportcard.com/report/github.com/vinci-proxy/vinci/mux)

Simple, versatile, general purpose HTTP multiplexer for [vinci](https://github.com/vinci-proxy/vinci) supporting multiple HTTP matching/filtering rules and easy composition capabilities.

This is more a convenient solution than very effecient one. If you're looking for great performance, simply pick another solution.

## Installation

```bash
go get -u gopkg.in/vinci-proxy/vinci.v0/mux
```

## API

See [godoc](https://godoc.org/github.com/vinci-proxy/mux) reference.

## Examples

#### Simple multiplexer

```go
package main

import (
  "fmt"
  "gopkg.in/vinci-proxy/mux.v0"
  "gopkg.in/vinci-proxy/vinci.v0"
  "net/http"
)

func main() {
  vs := vinci.NewServer(vinci.ServerOptions{Host: "localhost", Port: 3100})

  m := mux.New()
  m.If(mux.MatchMethod("GET", "POST"), mux.MatchPath("^/foo"))

  m.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
    w.Header().Set("Server", "vinci")
    h.ServeHTTP(w, r)
  })

  m.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
    w.Write([]byte("foo"))
  })

  vs.Use(m)
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", 3100)
  err := vs.Listen()
  if err != nil {
    fmt.Printf("Error: %s\n", err)
  }
}
```

#### Composition

```go
package main

import (
  "fmt"
  "gopkg.in/vinci-proxy/mux.v0"
  "gopkg.in/vinci-proxy/vinci.v0"
  "net/http"
)

func main() {
  vs := vinci.NewServer(vinci.ServerOptions{Host: "localhost", Port: 3100})

  // Create a custom multiplexer for /ip path
  ip := mux.If(mux.Path("^/ip"))
  ip.Use(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(r.RemoteAddr))
  })

  // Create a custom multiplexer for /headers path
  headers := mux.If(mux.Path("^/headers"))
  headers.Use(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(fmt.Errorf("Headers: %#v", r.Header).Error()))
  })

  // Creates the root multiplexer who host both multiplexers
  m := mux.New()
  m.If(mux.MatchMethod("GET"))
  m.Use(ip)
  m.Use(headers)

  // Register the multiplexer in the vinci
  vs.Use(m)
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", 3100)
  err := vs.Listen()
  if err != nil {
    fmt.Printf("Error: %s\n", err)
  }
}
```

#### Custom matcher function

```go
package main

import (
  "fmt"
  "gopkg.in/vinci-proxy/mux.v0"
  "gopkg.in/vinci-proxy/vinci.v0"
  "net/http"
)

func main() {
  vs := vinci.NewServer(vinci.ServerOptions{Host: "localhost", Port: 3100})

  m := mux.New()

  // Register a custom matcher function
  m.If(func(req *http.Request) bool {
    return req.Method == "GET" && req.RequestURI == "/foo"
  })

  m.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
    w.Header().Set("Server", "vinci")
    h.ServeHTTP(w, r)
  })

  m.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
    w.Write([]byte("foo"))
  })

  vs.Use(m)
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", 3100)
  err := vs.Listen()
  if err != nil {
    fmt.Printf("Error: %s\n", err)
  }
}
```

## License

MIT - Tomas Aparicio
