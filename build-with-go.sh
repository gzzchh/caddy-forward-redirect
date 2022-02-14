#!/bin/bash
#export GOFLAGS='-gcflags="all=-N -l"'
CGO_ENABLED=0 go build -o ../caddy-debug -a -ldflags '-extldflags "-static"' -gcflags="all=-N -l" .