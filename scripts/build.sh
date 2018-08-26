#!/usr/bin/env bash
dir="${0%/*}"

for os in "linux" "windows" "darwin"; do
	mkdir -p "$dir/../builds/$os"
	(cd "$dir/../" || exit; GOOS="$os" GOARCH=amd64 make build)
done
