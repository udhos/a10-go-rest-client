#!/bin/bash

build() {
        local path="$1"

        gofmt -s -w $path
        go tool fix $path
        go tool vet $path

	hash golint && golint $path

        CGO_ENABLED=0 go test $path
        CGO_ENABLED=0 go install $path
}

build ./a10go
build ./examples/a10list
build ./examples/a10server
build ./examples/a10sgroup

