#!/bin/sh

cd ../cmd/sample || exit

go build -ldflags="-s -w" -o ../../dist
