#!/usr/bin/env bash
export GO111MODULE=on
export CONSUL=false
rm -rf app
go build -o app
./app