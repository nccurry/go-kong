# go-kong #

[![GoDoc](https://godoc.org/github.com/spicyusername/go-kong/kong?status.svg)](https://godoc.org/github.com/spicyusername/go-kong/kong) 
[![Build Status](https://travis-ci.org/spicyusername/go-kong.svg?branch=master)](https://travis-ci.org/spicyusername/go-kong) 
[![Coverage Status](https://coveralls.io/repos/github/spicyusername/go-kong/badge.svg?branch=master)](https://coveralls.io/github/spicyusername/go-kong?branch=master) 
[![Go Report Card](https://goreportcard.com/badge/github.com/spicyusername/go-kong)](https://goreportcard.com/report/github.com/spicyusername/go-kong)  

go-kong is a Go client library for accessing Mashape's [Kong API](https://getkong.org/docs/0.9.x/admin-api/).

## Install ##

```bash
go get "github.com/spicyusername/go-kong/kong"
```

## Usage ##

Construct a new Kong client, then use the various services on the client to
access different parts of the Kong API. 

For example:

```go
package main

import (
        "log"
        "github.com/spicyusername/go-kong/kong"
)

func main() {
    // Create new client
    client, _ := kong.NewClient(nil, "http://localhost:8001/")
    
    // Get information about the 'backend' api
    api, _, _ := client.Apis.Get("backend")
    log.Println(api.ID)

    // Create a new api called 'middletier'
    mtApi := &kong.Api{Name: "middletier", RequestPath: "/mt/v0", UpstreamURL: "http://mt.my.org:8080"}
    _, err := client.Apis.Post(mtApi)
    if ok := err.(kong.ConflictError); ok {
        log.Printf("Endpoint with name %s already exists.", mtApi.Name)
    } else if err != nil {
        log.Fatal(err)
    }
    
    // Get all consumer objects
    consumers, _, _ := client.Consumers.GetAll(nil)
    for _, v := range consumers.Data {
        log.Println(v.Username)
    }

    // Apply ACL plugin to all all apis
    aclConfig := &kong.ACLConfig{Whitelist: []string{"users", "admins"}, Blacklist: []string{"blocked"}}
    plugin := &kong.Plugin{Name: "acl", Config: kong.ToMap(aclConfig)}
    _, err = client.Plugins.Post(plugin)
    if err != nil {
        log.Fatal(err)
    }
    
    // Add ACL group to consumer
    aclConsumerConfig := &kong.ACLConsumerConfig{Group: "users"}
    consumerName := "paul.atreides"
    _, err = client.Consumers.ACL.Configure(consumerName, aclConsumerConfig)
    if ok := err.(kong.NotFoundError); ok {
        log.Printf("Could not find consumer with name %s", consumerName)
    } else if err != nil {
        log.Fatal(err)
    }
}
```

## Apis ##
```go
// GET /apis
apis, resp, err := client.Apis.GetAll(nil)

// GET /apis?size=10&mapKey=myapi
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

// GET /plugins?size=15&mapKey=acl

// GET /plugins/4def15f5-0697-4956-a2b0-9ae079b686bb

// GET /plugins/enabled

// GET /plugins/schema/acl

// POST /plugins

// PATCH /plugins

// DELETE /plugins/4def15f5-0697-4956-a2b0-9ae079b686bb
```
