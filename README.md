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

    // (...)

    options := a10go.Options{Debug: true} // client options
    c := a10go.New(host, options) // create api client

    errLogin := c.Login(user, pass) // open session
    if errLogin != nil {
        fmt.Printf("login failure: %v", errLogin)
        return
    }
    vServers := c.VirtualServerList()

    c.Logout() // close session

See GoDoc: [http://godoc.org/github.com/udhos/a10-go-rest-client/a10go](http://godoc.org/github.com/udhos/a10-go-rest-client/a10go)

See example: [a10list](https://github.com/udhos/a10-go-rest-client/blob/master/examples/a10list/main.go)
