#!/usr/bin/env bash
dir="${0%/*}"
pids=()

for os in "linux" "windows" "darwin"; do
	mkdir -p "$dir/../builds/$os"
	(cd "$dir/../" || exit; GOOS="$os" GOARCH=amd64 make build) &
	pids[${os}]=$!
done

for pid in ${pids[*]}; do
	wait "$pid"
done
