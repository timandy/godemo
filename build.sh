#!/bin/bash

export PATH=$PATH:$GOPATH/bin
go build -toolexec="routiner go-agent" "$@"
