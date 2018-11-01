[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/a10-go-rest-client/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/udhos/a10-go-rest-client?status.svg)](http://godoc.org/github.com/udhos/a10-go-rest-client/a10go)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/a10-go-rest-client)](https://goreportcard.com/report/github.com/udhos/a10-go-rest-client)

# a10-go-rest-client
A10 golang rest client

# Build

    git clone https://github.com/udhos/a10-go-rest-client ;# clone outside of GOPATH
    cd a10-go-rest-client
    go install ./a10go
    go install ./examples/a10list

# Usage

    import "github.com/udhos/a10-go-rest-client/a10go"

    c := a10go.New(host, a10go.Options{Debug: true})
    errLogin := c.Login(user, pass)
    if errLogin != nil {
        fmt.Printf("login failure: %v", errLogin)
        return
    }
    vServers := c.VirtualServerList()

See GoDoc: [http://godoc.org/github.com/udhos/a10-go-rest-client/a10go](http://godoc.org/github.com/udhos/a10-go-rest-client/a10go)

See example: [a10list](https://github.com/udhos/a10-go-rest-client/blob/master/examples/a10list/main.go)
