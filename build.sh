#!/bin/bash

export GIT_HASH="$(git rev-parse HEAD)";
echo "Building git hash: ${GIT_HASH}";

export VERSION="0.0.1"

export BINARY="./bin/httpd"
export LDFLAGS="-s -w -buildid= -X main.version=${VERSION}"


go build  -x  -v  -trimpath \
    -tags="${TAGS}" \
    -ldflags="${LDFLAGS} \
    -X main.gitHash=${GIT_HASH}" \
    -o "${BINARY}"
