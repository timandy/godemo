@echo off

SET PATH=%PATH%;%GOPATH%/bin
go build -toolexec="routiner go-agent" %*
