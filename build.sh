#!/bin/bash

GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" kill_crashpad.go
