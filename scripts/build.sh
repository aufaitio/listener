#!/usr/bin/env bash

mkdir -p builds/linux builds/window builds/darwin

for os in "linux" "windows" "darwin"; do
    GOOS="$os" GOARCH=amd64 make build
done