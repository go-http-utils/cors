# cors
[![Build Status](https://travis-ci.org/go-http-utils/cors.svg?branch=master)](https://travis-ci.org/go-http-utils/cors)
[![Coverage Status](https://coveralls.io/repos/github/go-http-utils/cors/badge.svg?branch=master)](https://coveralls.io/github/go-http-utils/cors?branch=master)

CORS middleware for Go.

## Installation

```
go get -u github.com/go-http-utils/cors
```

## Documentation

API documentation can be found here: https://godoc.org/github.com/go-http-utils/cors

## Usage

```go
import (
  "github.com/go-http-utils/cors"
)
```

```go
mux := http.NewServeMux()
mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
  res.Write([]byte("Hello World"))
})

http.ListenAndServe(":8080", cors.Handler(mux))
```
