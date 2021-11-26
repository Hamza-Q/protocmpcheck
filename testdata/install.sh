#!/bin/sh
GOPATH=$PWD GO111MODULE=off go get github.com/stretchr/testify google.golang.org/grpc google.golang.org/grpc/examples/helloworld/helloworld
