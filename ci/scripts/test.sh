#!/bin/bash

set -eu

export GOPATH=$PWD/gopath
export PATH=$GOPATH/bin:$PATH

cd $GOPATH/src/github.com/cloudfoundry-community/cf-subway

go test $(go list ./... | grep -v /vendor/)
