#!/bin/bash

set -e -u -x

mkdir -p $GOPATH/src/github.com/legnoh
cp -r repo $GOPATH/src/github.com/legnoh/wlw-locate-kml
cd $GOPATH/src/github.com/legnoh/wlw-locate-kml

dep ensure
go run main.go
cd -

cp $GOPATH/src/github.com/legnoh/wlw-locate-kml/result-*.kml out/
