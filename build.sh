#!/usr/bin/env bash

#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build resource-update
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build resource-update
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build resource-update
