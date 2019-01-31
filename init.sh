#!/bin/bash
export PATH=`pwd`/bin:$PATH
go build -o bin/bld cmd/bld/build.go