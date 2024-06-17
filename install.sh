#!/bin/sh

BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
GIT_COMMIT="$(git rev-parse HEAD)"
VERSION="$(git describe --tags --abbrev=0 | tr -d '\n')"
echo "${BUILD_DATE} ${GIT_COMMIT} ${VERSION}"

go build -o dj -ldflags="-X 'github.com/ashish10alex/dj/internal/version.buildDate=${BUILD_DATE}' -X 'github.com/ashish10alex/dj/internal/version.gitCommit=${GIT_COMMIT}' -X 'github.com/ashish10alex/dj/internal/version.gitVersion=${VERSION}'" .

