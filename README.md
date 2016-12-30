# go-kong #

[![GoDoc](https://godoc.org/github.com/spicyusername/go-kong/kong?status.svg)](https://godoc.org/github.com/spicyusername/go-kong/kong) 
[![Build Status](https://travis-ci.org/spicyusername/go-kong.svg?branch=master)](https://travis-ci.org/spicyusername/go-kong) 
[![Coverage Status](https://coveralls.io/repos/github/spicyusername/go-kong/badge.svg?branch=master)](https://coveralls.io/github/spicyusername/go-kong?branch=master)  
[![Go Report Card](https://goreportcard.com/badge/github.com/spicyusername/go-kong)](https://goreportcard.com/report/github.com/spicyusername/go-kong)  

go-kong is a Go client library for accessing Mashape's [Kong API](https://getkong.org/docs/0.9.x/admin-api/).

## Install ##

```go
go get "github.com/spicyusername/go-kong/kong"
```

## Usage ##

Construct a new Kong client, then use the various services on the client to
access different parts of the Kong API. 

For example:

```go
package main

import (
        "github.com/spicyusername/go-kong/kong"
        "fmt"
)

func main() {
    client := kong.NewClient(nil, "http://localhost:8001/")
    
    // Get information about the 'backend' api
    api, _, _ := client.Apis.Get("backend")
    fmt.Printf("%+v", api)
    
    // Get all consumer objects
    consumers, _, _ := client.Consumers.GetAll(nil)
    fmt.Printf("%+v", consumers)
}

```

## Apis ##
```go
// GET /apis
apis, resp, err := client.Apis.GetAll(nil)

// GET /apis?size=10&name=myapi
opt := &ApisGetAllOptions{Size: 10, Name: "myapi"}
apis, resp, err := client.Apis.GetAll(opt)

// GET /apis/myapi
api, resp, err := client.Apis.Get("myapi")

// POST /apis
api := &Api{Name: "myapi", RequestPath: "/myapi", UpstreamURL: "http:myapi:8080"}
resp, err := client.Apis.Post(api)

// PATCH /apis/myapi
api := &Api{Name: "myapi", RequestPath: "/myapi", UpstreamURL: "http:myapi:8080"}
resp, err := client.Apis.Patch(api)

// DELETE /apis/myapi
resp, err := client.Apis.Delete("myapi")
```

## Consumers ##
```go
// GET /consumers
consumers, resp, err := client.Consumers.GetAll(nil)

// GET /consumers?size=10&custom_id=nickname
opt := &ConsumersGetAllOptions{Size: 10, CustomID: "nickname"}
consumers, resp, err := client.Consumers.GetAll(opt)

// GET /consumers/admin
consumer, resp, err := client.Consumers.Get("admin")

// POST /consumers
consumer := &Consumer{Username: "admin"}
resp, err := client.Consumers.Post(consumer)

// PATCH /consumers/admin
consumer := &Consumer{CustomID: "superuser"}
resp, err := client.Consumers.Patch(consumer)

// DELETE /consumers/admin
resp, err := client.Consumers.Delete("admin")
```

## Plugins ##

```go
// GET /plugins

// GET /plugins?size=15&name=acl

// GET /plugins/4def15f5-0697-4956-a2b0-9ae079b686bb

// GET /plugins/enabled

// GET /plugins/schema/acl

// POST /plugins

// PATCH /plugins

// DELETE /plugins/4def15f5-0697-4956-a2b0-9ae079b686bb
```
