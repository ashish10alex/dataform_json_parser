#!/bin/sh

VERSION=$(git describe --tags --abbrev=0 | tr -d '\n') GIT_COMMIT=$(git rev-parse HEAD) BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') goreleaser release --snapshot --clean
