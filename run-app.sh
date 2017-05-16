#!/bin/bash

export VERSION=v6.1.10
trap "exit 0" HUP INT QUIT TERM

go-cloud-debug -v -sourcecontext=source-context.json -appmodule=gke-info -appversion=${VERSION} -- ./gke-info &
./gke-info