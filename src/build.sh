#!/bin/bash

set -v on
export CGO_ENABLED=0

export SERVICE=push_service GOOS=linux GOARCH=amd64
go build -o "${SERVICE}-${GOOS}-${GOARCH}"
