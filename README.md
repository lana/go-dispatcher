# GO Dispatcher

[![Build Status](https://img.shields.io/travis/lana/go-dispatcher/master.svg?style=flat-square)](https://travis-ci.org/lana/go-dispatcher)
[![Codecov branch](https://img.shields.io/codecov/c/github/lana/go-dispatcher/master.svg?style=flat-square)](https://codecov.io/gh/lana/go-dispatcher)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/lana/go-dispatcher)
[![Go Report Card](https://goreportcard.com/badge/github.com/lana/go-dispatcher?style=flat-square)](https://goreportcard.com/report/github.com/lana/go-dispatcher)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://github.com/lana/go-dispatcher/blob/master/LICENSE)

A simple event dispatcher made in Go.

## Install

Use go get.
```sh
$ go get github.com/lana/go-dispatcher
```

Then import the package into your own code:
```
import "github.com/lana/go-dispacher"
```

## Usage
```go
package main

import (
	"context"
	"log"
	
	"github.com/lana/go-dispatcher"
)

type User struct {
	Name  string
	Email string
}

var UserWasCreated dispatcher.EventType = "user-was-created"

type UserCreated struct {
	User User
}

func (uc UserCreated) Type() dispatcher.EventType {
	return UserWasCreated
}

func (uc UserCreated) Data() interface{} {
	return uc.User
}

func OnUserCreated(ctx context.Context, e dispatcher.Event) {
	log.Printf("A user was created: %v", e.Data())
}

func main() {
	d := dispatcher.New()
	d.On(UserWasCreated, OnUserCreated)
	
	if err := d.Dispatch(context.Bachground(), e); err != nil {
		log.Fatalf("failed to dispatch event: %v", err)
	}
}
```

## License

This project is released under the MIT licence. See [LICENSE](https://github.com/lana/go-dispatcher/blob/master/LICENSE) for more details.
