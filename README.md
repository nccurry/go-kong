# go-kong #

[![GoDoc](https://godoc.org/github.com/spicyusername/go-kong/kong?status.svg)](https://godoc.org/github.com/spicyusername/go-kong/kong) 
[![Build Status](https://travis-ci.org/spicyusername/go-kong.svg?branch=master)](https://travis-ci.org/spicyusername/go-kong) 
[![Coverage Status](https://coveralls.io/repos/github/spicyusername/go-kong/badge.svg?branch=master)](https://coveralls.io/github/spicyusername/go-kong?branch=master)  

go-kong is a Go client library for accessing Mashape's [Kong API](https://getkong.org/docs/0.9.x/admin-api/).


## Usage ##

```go
import "github.com/spicyusername/go-kong/kong"
```

Construct a new Kong client, then use the various services on the client to
access different parts of the Kong API. 

For example:

```go
client := kong.NewClient(nil, "http://localhost:8001/")

// Get information about the 'backend' api
apis, _, err := client.Apis.Get("backend")
```