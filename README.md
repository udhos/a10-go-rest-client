[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/a10-go-rest-client/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/udhos/a10-go-rest-client?status.svg)](http://godoc.org/github.com/udhos/a10-go-rest-client/a10go)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/a10-go-rest-client)](https://goreportcard.com/report/github.com/udhos/a10-go-rest-client)

# a10-go-rest-client
A10 golang rest client

# Build

    git clone https://github.com/udhos/a10-go-rest-client ;# clone outside of GOPATH
    cd a10-go-rest-client
    go test ./a10go
    go install ./a10go
    go install ./examples/...

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

See [examples](https://github.com/udhos/a10-go-rest-client/tree/master/examples):

- [a10list](https://github.com/udhos/a10-go-rest-client/blob/master/examples/a10list/main.go)
- [a10server](https://github.com/udhos/a10-go-rest-client/blob/master/examples/a10server/main.go)
- [a10sgroup](https://github.com/udhos/a10-go-rest-client/blob/master/examples/a10sgroup/main.go)
- [a10vserver](https://github.com/udhos/a10-go-rest-client/blob/master/examples/a10vserver/main.go)

