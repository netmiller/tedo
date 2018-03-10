#!/bin/bash

#MAC OS
#export GOARCH="386"
#export GOOS="darwin"
#export CGO_ENABLED=1
#go build -o tedo -v

export GOARCH="amd64"
export GOOS="darwin"
export CGO_ENABLED=1
go build -o tedo.mac -v

#LINUX
export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=0
go build -o tedo.linux -v

#export GOARCH="386"
#export GOOS="linux"
#export CGO_ENABLED=0
#go build -o tedo.linux -v

#WINDOWS
export GOARCH="386"
export GOOS="windows"
export CGO_ENABLED=0
go build -o tedo32.exe -v

export GOARCH="amd64"
export GOOS="windows"
export CGO_ENABLED=0
go build -o tedo64.exe -v
